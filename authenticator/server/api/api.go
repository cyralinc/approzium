package api

import (
	"net/http"
	"strconv"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/cyralinc/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
)

func Start(logger *log.Logger, config config.Config) <-chan error {
	errChan := make(chan error)

	endpoints, err := loadEndpoints(logger, config)
	if err != nil {
		errChan <- err
		return errChan
	}

	serviceAddress := config.Host + ":" + strconv.Itoa(config.HTTPPort)
	logger.Infof("api listening for requests on %s", serviceAddress)
	go func() {
		if err := http.ListenAndServe(serviceAddress, endpoints); err != nil {
			errChan <- err
		}
	}()
	return errChan
}

func loadEndpoints(logger *log.Logger, config config.Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	prometheusHandler, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "approzium",
	})
	if err != nil {
		return nil, err
	}

	// Alphabetical by endpoint.
	mux.Handle("/v1/health", newHealthChecker(logger, config))
	mux.Handle("/v1/metrics/prometheus", prometheusHandler)

	return mux, nil
}
