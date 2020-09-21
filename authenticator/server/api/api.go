package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/cyralinc/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
)

// Start beings the authenticator's API server, and returns a function
// that can be called to begin graceful shutdown. Be sure to wait until
// the graceful shutdown function completes to end the program.
func Start(logger *log.Logger, config config.Config) (func(), error) {

	if err := loadEndpoints(logger, config); err != nil {
		return nil, err
	}

	serviceAddress := config.Host + ":" + strconv.Itoa(config.HTTPPort)
	server := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if config.DisableTLS {
		server.Addr = serviceAddress
		go func() {
			logger.Fatal(server.ListenAndServe())
		}()
		logger.Infof("api starting on http://%s", serviceAddress)
	} else {
		server.Addr = fmt.Sprintf(":%d", config.HTTPPort)
		crt, err := ioutil.ReadFile(config.PathToTLSCert)
		if err != nil {
			return nil, err
		}
		key, err := ioutil.ReadFile(config.PathToTLSKey)
		if err != nil {
			return nil, err
		}
		cert, err := tls.X509KeyPair(crt, key)
		if err != nil {
			return nil, err
		}
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   config.Host,
		}

		go func() {
			// ListenAndServeTLS always blocks indefinitely until it returns an error
			// describing how it stopped.
			if err := server.ListenAndServeTLS("", ""); err != nil {
				if err.Error() == http.ErrServerClosed.Error() {
					logger.Info(err)
				} else {
					logger.Fatal(err)
				}
			}
		}()
		logger.Infof("api starting on https://%s", serviceAddress)
	}
	return func() {
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Errorf("error shutting down gracefully: %s", err)
		}
	}, nil
}

func loadEndpoints(logger *log.Logger, config config.Config) error {
	prometheusHandler, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "approzium",
	})
	if err != nil {
		return err
	}

	// Alphabetical by endpoint.
	http.Handle("/v1/health", newHealthChecker(logger, config))
	http.Handle("/v1/metrics/prometheus", prometheusHandler)

	return nil
}
