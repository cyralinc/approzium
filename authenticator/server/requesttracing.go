package server

import (
	"context"
	"io"

	pb "github.com/cyralinc/approzium/authenticator/server/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

const tagRequestId = "request_id"

func newRequestTracer(wrapped pb.AuthenticatorServer) (pb.AuthenticatorServer, error) {
	tracer, closer, err := jaegerTracer()
	if err != nil {
		return nil, err
	}
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close() // TODO this supposed to be when the application ends
	return &requestTracer{
		wrapped: wrapped,
	}, nil
}

type requestTracer struct {
	wrapped pb.AuthenticatorServer
}

func (t *requestTracer) GetPGMD5Hash(ctx context.Context, req *pb.PGMD5HashRequest) (*pb.PGMD5Response, error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("GetPGMD5Hash")
	defer span.Finish()

	span.SetTag(tagRequestId, getRequestId(ctx))
	return t.wrapped.GetPGMD5Hash(ctx, req)
}

func (t *requestTracer) GetPGSHA256Hash(ctx context.Context, req *pb.PGSHA256HashRequest) (*pb.PGSHA256Response, error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("GetPGSHA256Hash")
	defer span.Finish()

	span.SetTag(tagRequestId, getRequestId(ctx))
	return t.wrapped.GetPGSHA256Hash(ctx, req)
}

func (t *requestTracer) GetMYSQLSHA1Hash(ctx context.Context, req *pb.MYSQLSHA1HashRequest) (*pb.MYSQLSHA1Response, error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("GetMYSQLSHA1Hash")
	defer span.Finish()

	span.SetTag(tagRequestId, getRequestId(ctx))
	return t.wrapped.GetMYSQLSHA1Hash(ctx, req)
}

func jaegerTracer() (opentracing.Tracer, io.Closer, error) {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: "approzium_authentication_server",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	return cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
}
