package main

import (
	"strings"

	"github.com/approzium/approzium/authenticator/server"
	"github.com/approzium/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	c, err := config.ParseConfig()
	if err != nil {
		log.Errorf("couldn't parse config: %s", err)
	}

	logLevel, err := log.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		log.Errorf("couldn't parse log level: %s", err)
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
		logger.Errorf("unsupported log format: %s", c.LogFormat)
	}

	if err := server.Start(logger, c); err != nil {
		logger.Errorf("authenticator ended due to %s", err)
	}
}
