package api

import (
	"net/http"
	"time"

	"encoding/json"

	"github.com/cyralinc/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
)

func newHealthChecker(logger *log.Logger, config config.Config) http.Handler {
	return &healthChecker{
		logger: logger,
		config: config,
	}
}

type healthChecker struct {
	logger *log.Logger
	config config.Config
}

const currentVersion = "0.2.0"

type healthResponse struct {
	ServerTime string `json:"server_time_utc"`
	Version    string `json:"version"`
}

// For most things that use a health check to determine if a service is up,
// like for AWS and Kubernetes for instance, a service is considered unhealthy
// if it returns anything other than a 200.
func (h *healthChecker) ServeHTTP(wr http.ResponseWriter, _ *http.Request) {
	h.logger.Debug("server is healthy")

	wr.Header().Set("Content-Type", "application/json")
	response := healthResponse{
		ServerTime: time.Now().UTC().Format(time.RFC3339),
		Version:    currentVersion,
	}

	h.logger.Debugf("healthResponse: %+v", response)

	json.NewEncoder(wr).Encode(response)
}
