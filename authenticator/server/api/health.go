package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/approzium/authenticator/server/config"
	protos "github.com/cyralinc/approzium/authenticator/server/protos"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	// Since _this_ endpoint is our mux API, we can conclude that it's working.
	// We need to also check that our gRPC server is up.
	grpcServerAddress := fmt.Sprintf("%s:%d", h.config.Host, h.config.GRPCPort)

	// It's fine to use an insecure connection here because we're
	// executing an unauthenticated health check that will return
	// nothing but a status code.
	conn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure()) // TODO would this work if TLS?
	if err != nil {
		h.logger.Warnf("unable to dial grpc due to %s", err)
		wr.WriteHeader(503)
		return
	}
	defer conn.Close()

	healthClient := protos.NewHealthClient(conn)
	resp, err := healthClient.Check(context.Background(), &protos.HealthCheckRequest{})
	if err != nil {
		h.logger.Warnf("unable to dial grpc due to %s", err)
		wr.WriteHeader(503)
		return
	}
	if resp.Status != protos.HealthCheckResponse_SERVING {
		h.logger.Warnf("unexpected response status %d", resp.Status)
		wr.WriteHeader(503)
		return
	}
	h.logger.Debug("server is healthy")
	wr.WriteHeader(200)
}
