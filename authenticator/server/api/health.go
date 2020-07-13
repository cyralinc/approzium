package api

import (
	"net/http"

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

// For most things that use a health check to determine if a service is up,
// like for AWS and Kubernetes for instance, a service is considered unhealthy
// if it returns anything other than a 200.
func (h *healthChecker) ServeHTTP(wr http.ResponseWriter, _ *http.Request) {
	h.logger.Debug("server is healthy")
	wr.WriteHeader(200)
}
