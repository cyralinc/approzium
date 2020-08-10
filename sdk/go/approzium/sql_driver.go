package approzium

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"

	pb "github.com/cyralinc/approzium/authenticator/server/protos"
	"github.com/cyralinc/approzium/sdk/go/approzium/identity"
	"github.com/cyralinc/pq"
	_ "github.com/cyralinc/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultPostgresPort = "5432"
	postgresUrlPrefix   = "postgres://"
)

var passwordIncludedErr = errors.New("approzium is for passwordless authentication and uses your identity as your password, " +
	"please remove the password field from your connection string")

// Examples of grpcAddr:
// 		- authenticator:6001
// 		- localhost:6001
// 		- somewhere:6001
// HTTPS is used unless TLS is disabled.
func NewAuthClient(grpcAddress string, config *Config) (*AuthClient, error) {
	if err := config.parse(); err != nil {
		return nil, err
	}

	identityHandler, err := identity.NewHandler(config.Logger, config.RoleArnToAssume)
	if err != nil {
		return nil, err
	}
	return &AuthClient{
		grpcAddress:     grpcAddress,
		config:          config,
		identityHandler: identityHandler,
	}, nil
}

type AuthClient struct {
	grpcAddress string
	config      *Config

	// This is used for preventing races, because we overwrite the pq library's global hashing
	// func each time Open is called.
	hashFuncLock      sync.Mutex
	openHasBeenCalled bool

	// This is used for caching identity for an appropriate period of time.
	identityHandler *identity.Handler
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a
// database name and connection information.
//
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once.
func (a *AuthClient) Open(driverName, dataSourceName string) (*sql.DB, error) {
	switch driverName {
	case "postgres":
		return a.openPostgres(driverName, dataSourceName)
	default:
		return nil, fmt.Errorf("%s is not supported", driverName)
	}
}

func (a *AuthClient) openPostgres(driverName, dataSourceName string) (*sql.DB, error) {
	dataSourceName, err := addPlaceholderPassword(dataSourceName)
	if err != nil {
		return nil, err
	}

	dbHost, dbPort, err := parseDSN(a.config.Logger, dataSourceName)
	if err != nil {
		return nil, err
	}

	/*
		When we update the GetMD5Hash function below, it'll be used later every time the database is
		queried. So we can think of it usually as, you create an AuthClient once, call Open once, then
		call multiple queries.

		However, in Go, it's trivial to parallelize your code to speed it up. So someone may create an
		AuthClient, then break into multiple goroutines, call Open on multiple goroutines, then commence
		with queries.

		However, if they had DIFFERENT hosts, ports, or AuthClient configs on each thread, we would have a
		problem. That's because the GetMD5Hash function is set at a GLOBAL level in the pq library, and called
		in each query. If there were differing Open calls on each goroutine, and then each goroutine had subsequent
		queries relying upon it, the call would constantly change underfoot of the queries in an illogical
		way.

		Here, we're preventing that race condition by holding a lock while we check that Open has only been called
		once, and while we update the global fetchHashFromApproziumAuthenticator function.
	*/
	a.hashFuncLock.Lock()
	defer a.hashFuncLock.Unlock()
	if a.openHasBeenCalled {
		return nil, errors.New("to prevent races, Open can only be called once, please instantiate a new authenticator" +
			"or only call Open once per AuthClient")
	}
	pq.GetMD5Hash = a.fetchHashFromApproziumAuthenticator(dbHost, dbPort)
	a.openHasBeenCalled = true

	return sql.Open(driverName, dataSourceName)
}

func (a *AuthClient) fetchHashFromApproziumAuthenticator(dbHost, dbPort string) func(user, password, salt string) (string, error) {
	return func(user, password, salt string) (string, error) {
		conn, err := a.grpcConnection()
		if err != nil {
			return "", err
		}
		defer func() {
			if err := conn.Close(); err != nil {
				// Failing to close connections may lead to a memory or connection leak,
				// so we should call attention to it.
				a.config.Logger.Warnf("unable to close GRPC connection due to %s", err)
			}
		}()
		client := pb.NewAuthenticatorClient(conn)

		// We call retrieve just before every call because the proof should not be cached by us -
		// the underlying identity handler deals with caching where appropriate.
		proof := a.identityHandler.Retrieve()
		resp, err := client.GetPGMD5Hash(context.Background(), &pb.PGMD5HashRequest{
			PwdRequest: &pb.PasswordRequest{
				ClientLanguage: proof.ClientLang,
				Dbhost:         dbHost,
				Dbport:         dbPort,
				Dbuser:         user,
				Aws:            proof.AwsAuth,
			},
			Salt: []byte(salt),
		})
		if err != nil {
			return "", err
		}
		return resp.Hash, nil
	}
}

func (a *AuthClient) grpcConnection() (*grpc.ClientConn, error) {
	if a.config.DisableTLS {
		return grpc.Dial(a.grpcAddress, grpc.WithInsecure())
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: a.config.InsecureSkipVerify,
	}

	if a.config.PathToClientCert != "" {
		clientCert, err := tls.LoadX509KeyPair(a.config.PathToClientCert, a.config.PathToClientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	if a.config.PathToTrustedCACerts != "" {
		trustedCACerts, err := ioutil.ReadFile(a.config.PathToTrustedCACerts)
		if err != nil {
			return nil, err
		}
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		if !pool.AppendCertsFromPEM(trustedCACerts) {
			return nil, fmt.Errorf("credentials: failed to append certificates at %s with body %s, "+
				"please check that they are valid CA certificates", a.config.PathToTrustedCACerts, trustedCACerts)
		}
		tlsConfig.RootCAs = pool
	}

	return grpc.Dial(a.grpcAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}

type Config struct {
	// Logger is optional. It's available for you to set in case you'd like to
	// customize it. If not set, it defaults to INFO level and text output.
	Logger *log.Logger

	// Set to true to disable. TLS is enabled by default.
	DisableTLS bool

	// Set to true to skip verifying the chain of trust on the server's
	// certificate.
	InsecureSkipVerify bool

	// This client's certificate, used for proving its identity, and used by
	// the caller to encrypt communication with its public key.
	PathToClientCert string

	// This client's key, used for decrypting incoming communication that was
	// encrypted by callers using the client cert's public key.
	PathToClientKey string

	// The path to the root certificate(s) that must have issued the identity
	// certificate used by Approzium's authentication server.
	PathToTrustedCACerts string

	// RoleArnToAssume is an optional field. It's useful for both testing, and
	// in environments like AWS Lambda where you'd to pull an ARN from the
	// enclosing environment to assume its identity.
	RoleArnToAssume string
}

func (c *Config) parse() error {
	if c.Logger == nil {
		c.Logger = log.New()
		c.Logger.SetLevel(log.InfoLevel)
		c.Logger.SetFormatter(&log.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		})
	}
	if !c.DisableTLS {
		if c.PathToClientCert == "" && c.PathToTrustedCACerts == "" {
			return errors.New("if TLS isn't disabled, the path to the TLS client certificate " +
				"or the path to trusted CA certs must be provided")
		}
		if c.PathToClientCert != "" && c.PathToClientKey == "" {
			return errors.New("if a client cert is supplied, a client key must be supplied as well")
		}
	}
	return nil
}

// addPlaceholderPassword ensures the user hasn't provided a password
// (because only the Approzium authentication server should have it),
// and then adds a placeholder password so lib/pq won't trip from not
// having anything supplied.
func addPlaceholderPassword(dataSourceName string) (string, error) {
	if !strings.HasPrefix(dataSourceName, postgresUrlPrefix) {
		// We received a string like:
		// user=postgres dbname=postgres host=localhost port=5432 sslmode=disable
		if strings.Contains(strings.ToLower(dataSourceName), "password") {
			return "", passwordIncludedErr
		}

		// Just add a password=unknown field to the end and return.
		return dataSourceName + " password=unknown", nil
	}

	// Convert strings like:
	//		"postgres://pqgotest:@localhost/pqgotest?sslmode=verify-full"
	// to:
	//		"postgres://pqgotest:unknown@localhost/pqgotest?sslmode=verify-full"
	u, err := url.Parse(dataSourceName)
	if err != nil {
		return "", err
	}
	if password, _ := u.User.Password(); password != "" {
		return "", passwordIncludedErr
	}

	fields := strings.Split(dataSourceName, "@")
	if len(fields) != 2 {
		return "", fmt.Errorf(`expected connection string like 'postgres://pqgotest:@localhost/pqgotest?sslmode=verify-full' but received %q`, dataSourceName)
	}
	return fields[0] + "unknown@" + fields[1], nil
}

func parseDSN(logger *log.Logger, dataSourceName string) (dbHost, dbPort string, err error) {
	if strings.HasPrefix(dataSourceName, postgresUrlPrefix) {
		u, err := url.Parse(dataSourceName)
		if err != nil {
			return "", "", err
		}
		// If a host and port were sent, this will initially come through as "localhost:1234"
		hostFields := strings.Split(u.Host, ":")
		dbHost = hostFields[0]
		dbPort = u.Port()
	} else {
		// Extract the host and port from a string like:
		// 		"user=postgres password=mysecretpassword dbname=postgres host=localhost port=5432 sslmode=disable"
		fields := strings.Split(dataSourceName, " ")
		for _, field := range fields {
			if dbHost != "" && dbPort != "" {
				break
			}
			kv := strings.Split(field, "=")
			if len(kv) != 2 {
				return "", "", fmt.Errorf("expected one = per group, but received %s", field)
			}
			key := kv[0]
			val := kv[1]
			if key == "host" {
				dbHost = val
				continue
			}
			if key == "port" {
				dbPort = val
				continue
			}
		}
	}

	if dbHost == "" {
		return "", "", fmt.Errorf("unable to parse host from %s", dataSourceName)
	}
	if dbPort == "" {
		logger.Warnf("unable to parse port from %s, defaulting to %s", dataSourceName, defaultPostgresPort)
		dbPort = defaultPostgresPort
	}
	return dbHost, dbPort, nil
}
