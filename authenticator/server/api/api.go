package api

import (
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	log "github.com/sirupsen/logrus"
)

func Start(logger *log.Logger, host, port string) <-chan error {
	errChan := make(chan error)

	endpoints, err := loadEndpoints()
	if err != nil {
		errChan <- err
		return errChan
	}

	serviceAddress := host + ":" + port
	logger.Infof("api listening for requests on %s", serviceAddress)
	go func() {
		if err := http.ListenAndServe(serviceAddress, endpoints); err != nil {
			errChan <- err
		}
	}()
	return errChan
}

func loadEndpoints() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	prometheusHandler, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "approzium",
	})
	if err != nil {
		return nil, err
	}
	mux.Handle("/v1/metrics/prometheus", prometheusHandler)

	return mux, nil
}
