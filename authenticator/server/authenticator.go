package server

/*

This server is stateless - it doesn't currently cache anything, perform
any writes, or have knowledge of other Approzium clusters. Because of this,
it can be highly available simply by running multiple instances. Please
do not add code that caches state unless we are planning to change to
a stateful, clustered design. Thanks!

*/

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/approzium/approzium/authenticator/server/api"
	"github.com/approzium/approzium/authenticator/server/config"
	"github.com/approzium/approzium/authenticator/server/credmgrs"
	"github.com/approzium/approzium/authenticator/server/identity"
	pb "github.com/approzium/approzium/authenticator/server/protos"
	"github.com/aws/aws-sdk-go/aws/arn"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The initial choice for a max is based on
// https://www.postgresql.org/docs/8.3/pgcrypto.html and
// https://tools.ietf.org/html/rfc7677.
// We want to allow enough iterations to be secure, but
// not so many that the iterations could be used to effectively
// DOS us by sending us looping for a long amount of time.
// The RFC recommends at least 15,000 iterations, so we just
// allow up to 10 times as much in case folks are being extra
// secure. We are open to making this higher or lower based on
// community feedback.
var maxIterations = uint32(15000 * 10)

// Start begins a GRPC server, and an API server. It hangs indefinitely until
// an error is returned from either, terminating the application. Both servers
// respond to CTRL+C shutdowns.
func Start(logger *log.Logger, config config.Config) error {
	apiErrChan := api.Start(logger, config)

	svr, err := buildServer(logger, config)
	if err != nil {
		return err
	}
	grpcErrChan := startGrpc(logger, config, svr)

	select {
	case err = <-apiErrChan:
	case err = <-grpcErrChan:
	}
	return err
}

func buildServer(logger *log.Logger, config config.Config) (pb.AuthenticatorServer, error) {
	// Calls pass through the following layers during handling.
	// 	- First, a layer that captures request metrics.
	//	- Next, a layer that adds a request ID, creates a request logger, and logs
	//		all inbound and outbound requests.
	//	- Lastly, this layer, the authenticator, that handles logic.
	authenticator, err := newAuthenticator(logger, config)
	if err != nil {
		return nil, err
	}

	svr, err := newRequestMetrics(newRequestLogger(logger, config.LogRaw, authenticator))
	if err != nil {
		return nil, err
	}
	return svr, nil
}

func startGrpc(logger *log.Logger, config config.Config, authenticatorServer pb.AuthenticatorServer) <-chan error {
	errChan := make(chan error)

	serviceAddress := fmt.Sprintf("%s:%d", config.Host, config.GRPCPort)
	lis, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		errChan <- err
		return errChan
	}

	logger.Infof("grpc listening for requests on %s", serviceAddress)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticatorServer(grpcServer, authenticatorServer)
	pb.RegisterHealthServer(grpcServer, newHealthServer())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- err
		}
	}()
	return errChan
}

func newAuthenticator(logger *log.Logger, config config.Config) (pb.AuthenticatorServer, error) {
	credMgr, err := credmgrs.RetrieveConfigured(logger, config.VaultTokenPath)
	if err != nil {
		return nil, err
	}
	identityVerifier, err := identity.NewVerifier()
	if err != nil {
		return nil, err
	}
	return &authenticator{
		logger:           logger,
		credMgr:          credMgr,
		identityVerifier: identityVerifier,
	}, nil
}

type authenticator struct {
	// The authenticator's logger lacks context about requests being made
	// and should not be used within code that's part of executing a request.
	logger           *log.Logger
	credMgr          credmgrs.CredentialManager
	identityVerifier identity.Verifier
}

func (a *authenticator) getPassword(reqLogger *log.Entry, req *pb.PasswordRequest) (string, error) {
	// Currently, only AWS identity is supported
	awsIdentity := req.GetAws()
	if awsIdentity == nil {
		return "", fmt.Errorf("AWS auth info is required")
	}
	// To expedite handling the request, let's verify the caller's identity at the same
	// time as getting the password.
	verifiedIdentityChan := make(chan *identity.Verified, 1)
	verificationErrChan := make(chan error, 1)
	go func() {
		proof := &identity.Proof{
			ClientLang: req.ClientLanguage,
			AwsAuth:    req.Aws,
		}
		verifiedIdentity, err := a.identityVerifier.Get(reqLogger, proof)
		if err != nil {
			verificationErrChan <- status.Errorf(codes.Unauthenticated, err.Error())
			return
		}
		verifiedIdentityChan <- verifiedIdentity
	}()

	claimedIamArn := awsIdentity.ClaimedIamArn
	dbHost := req.GetDbhost()
	dbPort := req.GetDbport()
	dbUser := req.GetDbuser()

	databaseArn, err := toDatabaseARN(reqLogger, claimedIamArn)
	if err != nil {
		return "", err
	}
	password, err := a.getCreds(reqLogger, credmgrs.DBKey{
		IAMArn: databaseArn,
		DBHost: dbHost,
		DBPort: dbPort,
		DBUser: dbUser,
	})
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Make sure the arn they claimed they had to get the creds was their actual arn.
	select {
	case verifiedIdentity := <-verifiedIdentityChan:
		match, err := a.identityVerifier.Matches(reqLogger, claimedIamArn, verifiedIdentity)
		if err != nil {
			return "", err
		}
		if !match {
			return "", status.Errorf(codes.Unauthenticated, fmt.Sprintf("claimed IAM arn %s did not match actual IAM arn of %+v", claimedIamArn, verifiedIdentity))
		}
	case err = <-verificationErrChan:
		return "", err
	}
	return password, nil
}

func (a *authenticator) GetPGMD5Hash(ctx context.Context, req *pb.PGMD5HashRequest) (*pb.PGMD5Response, error) {
	// Return early if we didn't get a valid salt.
	salt := req.GetSalt()
	if len(salt) != 4 {
		msg := fmt.Sprintf("expected salt to be 4 bytes long, but got %d bytes", len(salt))
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	reqLogger := getRequestLogger(ctx)
	password, err := a.getPassword(reqLogger, req.GetPwdRequest())
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	dbUser := req.GetPwdRequest().GetDbuser()
	// Everything checked out.
	hash, err := computePGMD5Hash(dbUser, password, salt)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &pb.PGMD5Response{Hash: hash}, nil
}

func (a *authenticator) GetPGSHA256Hash(ctx context.Context, req *pb.PGSHA256HashRequest) (*pb.PGSHA256Response, error) {
	// Return early if we didn't get a valid auth message or salt.
	authMsg := req.GetAuthenticationMsg()
	if len(authMsg) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("authentication message not provided"))
	}

	salt := req.GetSalt()
	if len(salt) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("salt not provided"))
	}

	iterations := req.GetIterations()
	if iterations > maxIterations {
		// Using a very high number of iterations could cause us to loop and a lot of
		// those requests could quickly take us down, like a DOS attack.
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("iterations too high, received %d but maximum is %d", iterations, maxIterations))
	}

	reqLogger := getRequestLogger(ctx)
	password, err := a.getPassword(reqLogger, req.GetPwdRequest())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	saltedPass, err := computePGSHA256SaltedPass(password, salt, int(iterations))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not compute hash %s", err))
	}

	cproof, err := computePGSHA256Cproof(saltedPass, authMsg)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	sproof := computePGSHA256Sproof(saltedPass, authMsg)
	return &pb.PGSHA256Response{Cproof: cproof, Sproof: sproof}, nil
}

func (a *authenticator) getCreds(reqLogger *log.Entry, identity credmgrs.DBKey) (string, error) {
	creds, err := a.credMgr.Password(reqLogger, identity)
	if err != nil {
		return "", fmt.Errorf("password not found for identity %s due to %s, using %s", identity, err, a.credMgr.Name())
	}
	return creds, nil
}

// toDatabaseARN either uses the original ARN to check the database
// for a password, or if it's an assumed role ARN, converts it to a
// role ARN before looking.
func toDatabaseARN(logger *log.Entry, fullIAMArn string) (string, error) {
	parsedArn, err := arn.Parse(fullIAMArn)
	if err != nil {
		return "", err
	}
	logger.Debugf("received login attempt from %+v", parsedArn)
	if !strings.HasPrefix(parsedArn.Resource, "assumed-role") {
		// This is a regular arn, so we should return it as-is for use in accessing
		// database credentials.
		return fullIAMArn, nil
	}
	// Convert assumed role arns to role arns.
	fields := strings.Split(parsedArn.Resource, "/")
	if len(fields) < 2 || len(fields) > 3 {
		return "", fmt.Errorf("unexpected assume role arn format: %s", fullIAMArn)
	}
	return fmt.Sprintf("arn:%s:iam::%s:role/%s", parsedArn.Partition, parsedArn.AccountID, fields[1]), nil
}

func computeMD5(s string, salt []byte) (string, error) {
	hasher := md5.New()
	if _, err := io.WriteString(hasher, s); err != nil {
		return "", err
	}
	hasher.Write(salt)
	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes), nil
}

func computePGMD5Hash(user, password string, salt []byte) (string, error) {
	firstHash, err := computeMD5(password, []byte(user))
	if err != nil {
		return "", err
	}
	secondHash, err := computeMD5(firstHash, salt)
	if err != nil {
		return "", err
	}
	return secondHash, nil
}

func computePGSHA256SaltedPass(password string, salt string, iterations int) ([]byte, error) {
	s, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, fmt.Errorf("Bad salt %s", err)
	}
	dk := pbkdf2.Key([]byte(password), s, iterations, 32, sha256.New)
	return dk, nil
}

// assumes a and b are of the same length
func xorBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("cannot xor slices of unequal lengths, received %d and %d", len(a), len(b))
	}
	buf := make([]byte, len(a))

	for i := range a {
		buf[i] = a[i] ^ b[i]
	}

	return buf, nil
}

// SCRAM reference: https://en.wikipedia.org/wiki/Salted_Challenge_Response_Authentication_Mechanism
func computePGSHA256Cproof(spassword []byte, authMsg string) (string, error) {
	mac := hmac.New(sha256.New, spassword)
	mac.Write([]byte("Client Key"))
	ckey := mac.Sum(nil)
	ckeyHash := sha256.Sum256(ckey)
	cproofHmac := hmac.New(sha256.New, ckeyHash[:])
	cproofHmac.Write([]byte(authMsg))
	cproof, err := xorBytes(cproofHmac.Sum(nil), ckey)
	if err != nil {
		return "", err
	}
	cproof64 := base64.StdEncoding.EncodeToString(cproof)
	return cproof64, nil
}

func computePGSHA256Sproof(spassword []byte, authMsg string) string {
	mac := hmac.New(sha256.New, spassword)
	mac.Write([]byte("Server Key"))
	skey := mac.Sum(nil)
	sproofHmac := hmac.New(sha256.New, skey)
	sproofHmac.Write([]byte(authMsg))
	sproof := sproofHmac.Sum(nil)
	sproof64 := base64.StdEncoding.EncodeToString(sproof)
	return sproof64
}
