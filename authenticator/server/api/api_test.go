package api

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/cyralinc/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
)

func TestAPI(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	i := strings.Index(wd, "approzium")
	if i == -1 {
		// This test is being run from outside the Approzium home directory.
		t.Skip("skipping because unable to locate the approzium home directory")
	}
	pathToTestDir := wd[:i] + "approzium/authenticator/server/testing"

	gracefulShutdown, err := Start(log.New(), config.Config{
		Host:          "127.0.0.1",
		HTTPPort:      6010,
		DisableTLS:    false,
		PathToTLSCert: pathToTestDir + "/approzium.pem",
		PathToTLSKey:  pathToTestDir + "/approzium.key",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer gracefulShutdown()

	// For this test, forgive that the approzium.pem cert is self-signed and is not for localhost.
	def := http.DefaultTransport
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	defer func() {
		http.DefaultTransport = def
	}()

	resp, err := http.Get("https://localhost:6010/v1/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("%d: %s", resp.StatusCode, b)
	}
}
