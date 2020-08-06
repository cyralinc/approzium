package testing

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/cyralinc/approzium/sdk/go/approzium"
	log "github.com/sirupsen/logrus"
)

func TestPostgresMD5(t *testing.T) {
	disableTLS := false
	if raw := os.Getenv("APPROZIUM_DISABLE_TLS"); raw != "" {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			t.Fatalf("couldn't parse APPROZIUM_DISABLE_TLS %s as a bool", raw)
		}
		disableTLS = b
	}
	authClient, err := approzium.NewAuthClient("authenticatortls:6001", &approzium.Config{
		Logger:               log.New(),
		DisableTLS:           disableTLS,
		PathToTrustedCACerts: os.Getenv("TEST_CERT_DIR") + "/approzium.pem",
		PathToClientCert:     os.Getenv("TEST_CERT_DIR") + "/client.pem",
		PathToClientKey:      os.Getenv("TEST_CERT_DIR") + "/client.key",
		RoleArnToAssume:      os.Getenv("TEST_ASSUMABLE_ARN"),

		// Because the server cert was issued by a self-signed certificate, it won't
		// pass verification. However, we still communicate via encrypted communication.
		InsecureSkipVerify: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceName := fmt.Sprintf("user=%s dbname=%s host=dbmd5 port=%s sslmode=require",
		os.Getenv("PSYCOPG2_TESTDB_USER"),
		os.Getenv("PSYCOPG2_TESTDB"),
		os.Getenv("PSYCOPG2_TESTDB_PORT"),
	)
	db, err := authClient.Open("postgres", dataSourceName)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("received nothing")
	}
}
