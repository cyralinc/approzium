package api

import (
	"testing"

	"github.com/cyralinc/approzium/authenticator/server/config"
	testtools "github.com/cyralinc/approzium/authenticator/server/testing"
)

func TestHealthChecker(t *testing.T) {
	checker := newHealthChecker(testtools.TestLogger(), config.Config{
		Host:     "127.0.0.1",
		Port: 6001,
	})
	testWriter := &testtools.TestResponseWriter{}
	checker.ServeHTTP(testWriter, nil)

	if testWriter.LastStatusCodeReceived != 200 {
		t.Fatalf("expected 200 but received %d", testWriter.LastStatusCodeReceived)
	}
}
