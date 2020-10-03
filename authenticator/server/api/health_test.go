package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyralinc/approzium/authenticator/server/config"
	testtools "github.com/cyralinc/approzium/authenticator/server/testing"
)

func TestHealthChecker(t *testing.T) {
	checker := newHealthChecker(testtools.TestLogger(), config.Config{
		Host:     "127.0.0.1",
		GRPCPort: 6001,
	})

	req, err := http.NewRequest("GET", "/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	checker.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Fatalf("expected %v but received %d", http.StatusOK, status)
	}

	var healthStatus healthResponse

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&healthStatus)
	if err != nil {
		t.Fatal(err)
	}
}
