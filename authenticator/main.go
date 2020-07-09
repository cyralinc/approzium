package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/approzium/approzium/authenticator/server"
	"github.com/approzium/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
)

var devMode bool

func init() {
	flag.BoolVar(&devMode, "dev", false, "run in dev mode")
}

func main() {
	flag.Parse()

	var c config.Config
	var err error
	if devMode {
		// dev mode uses the file back-end
		os.Unsetenv("VAULT_ADDR")
		c = config.Config{
			Host:       "127.0.0.1",
			HTTPPort:   6000,
			GRPCPort:   6001,
			DisableTLS: true,
			LogLevel:   "debug",
			LogFormat:  "text",
			LogRaw:     false,
		}
	} else {
		c, err = config.ParseConfig()
	}
	if err != nil {
		log.Errorf("couldn't parse config: %s", err)
		return
	}

	logger, err := buildApplicationLogger(c)
	if err != nil {
		log.Error(err)
		return
	}

	if err := server.Start(logger, c); err != nil {
		logger.Errorf("authenticator ended due to %s", err)
		return
	}
	logger.Info("all ports up and ready to serve traffic")

	// Wait for a shutdown signal.
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
}

func buildApplicationLogger(c config.Config) (*log.Logger, error) {
	logLevel, err := log.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		return nil, err
	}
	logger := log.New()
	logger.Level = logLevel

	switch strings.ToLower(c.LogFormat) {
	case "text":
		logger.SetFormatter(&log.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		})
	case "json":
		logger.SetFormatter(&log.JSONFormatter{})
	default:
		return nil, fmt.Errorf("unsupported log format: %s", c.LogFormat)
	}
	return logger, nil
}
