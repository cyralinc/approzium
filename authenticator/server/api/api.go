package api

import (
	"crypto/tls"
	"io/ioutil"
    "net"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/cyralinc/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
)

func Start(logger *log.Logger, listener net.Listener, config config.Config) error {

	if err := loadEndpoints(logger, config); err != nil {
		return err
	}

	server := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if config.DisableTLS {
		go func() {
			logger.Fatal(server.Serve(listener))
		}()
		logger.Infof("api server starting without TLS")
	} else {
		crt, err := ioutil.ReadFile(config.PathToTLSCert)
		if err != nil {
			return err
		}
		key, err := ioutil.ReadFile(config.PathToTLSKey)
		if err != nil {
			return err
		}
		cert, err := tls.X509KeyPair(crt, key)
		if err != nil {
			return err
		}
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   config.Host,
		}

		go func() {
			logger.Fatal(server.ServeTLS(listener, "", ""))
		}()
		logger.Infof("api server starting with TLS")
	}
	return nil
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
