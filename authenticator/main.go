package main

import (
	"strings"

	"github.com/approzium/approzium/authenticator/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := server.ParseConfig()
	if err != nil {
		log.Errorf("couldn't parse config: %s", err)
	}

	logLevel, err := log.ParseLevel(strings.ToLower(config.LogLevel))
	if err != nil {
		log.Errorf("couldn't parse log level: %s", err)
	}
	logger := log.New()
	logger.Level = logLevel

	switch strings.ToLower(config.LogFormat) {
	case "text":
		logger.SetFormatter(&log.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		})
	case "json":
		logger.SetFormatter(&log.JSONFormatter{})
	default:
		logger.Errorf("unsupported log format: %s", config.LogFormat)
	}

	if err := server.Start(logger, config); err != nil {
		logger.Errorf("authenticator ended due to %s", err)
	}
}
