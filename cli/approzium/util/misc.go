package util

import log "github.com/sirupsen/logrus"

// We reuse approzium server functions that output logs. We don't want
// to see them in the CLI, so we populate those functions with a silent
// logger.
var SilentLogger = func() *log.Logger {
	l := log.New()
	l.Level = log.PanicLevel
	return l
}()
