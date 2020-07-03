package testing

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func TestLogger() *log.Logger {
	logger := log.New()
	logger.Level = log.FatalLevel
	return logger
}

type TestResponseWriter struct {
	HeaderToReturn http.Header

	LastWriteInputReceived []byte
	WriteIntToReturn       int
	WriteErrToReturn       error

	LastStatusCodeReceived int
}

func (t *TestResponseWriter) Header() http.Header {
	return t.HeaderToReturn
}

func (t *TestResponseWriter) Write(input []byte) (int, error) {
	t.LastWriteInputReceived = input
	return t.WriteIntToReturn, t.WriteErrToReturn
}

func (t *TestResponseWriter) WriteHeader(statusCode int) {
	t.LastStatusCodeReceived = statusCode
}
