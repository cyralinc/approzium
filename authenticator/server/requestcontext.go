package server

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

type EmbellishedContext struct {
	// GRPC is the context given to us at a per-request level by the
	// GRPC framework.
	GRPC context.Context

	// RequestId is a randomly generated ID given by the logging handler.
	RequestId string

	// RequestLogger is a logger that also includes the request ID on ever
	// line. It should be used wherever possible.
	RequestLogger *log.Entry

	// During each call, a trace is started that is used for the lifetime
	// of the call. When we call out to _other_ systems, we should also
	// have a unique trace for that call that includes this tracing context.
	// This ties the sub-call to the overall call.
	TracingContext opentracing.SpanContext
}

func (c *EmbellishedContext) Deadline() (deadline time.Time, ok bool) {
	return c.GRPC.Deadline()
}

func (c *EmbellishedContext) Done() <-chan struct{} {
	return c.GRPC.Done()
}

func (c *EmbellishedContext) Err() error {
	return c.GRPC.Err()
}

func (c *EmbellishedContext) Value(key interface{}) interface{} {
	return c.GRPC.Value(key)
}
