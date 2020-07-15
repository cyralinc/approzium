package server

import (
	"context"
	"io"

	pb "github.com/cyralinc/approzium/authenticator/server/protos"
	"github.com/cyralinc/approzium/authenticator/server/tracing"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// https://github.com/jaegertracing/jaeger-client-go
const tagRequestId = "request_id"

func newRequestTracer(logger *log.Logger, wrapped pb.AuthenticatorServer) (pb.AuthenticatorServer, io.Closer, error) {
	tracer, gracefulShutdown, err := jaegerTracer(logger)
	if err != nil {
		return nil, nil, err
	}
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	return &requestTracer{
		wrapped: wrapped,
	}, gracefulShutdown, nil
}

type requestTracer struct {
	wrapped pb.AuthenticatorServer
}

func (t *requestTracer) GetPGMD5Hash(ctx context.Context, req *pb.PGMD5HashRequest) (*pb.PGMD5Response, error) {
	span := tracing.Tracer.StartSpan("GetPGMD5Hash")
	defer span.Finish()

	span.SetTag(tagRequestId, getRequestId(ctx))
	return t.wrapped.GetPGMD5Hash(ctx, req)
}

func (t *requestTracer) GetPGSHA256Hash(ctx context.Context, req *pb.PGSHA256HashRequest) (*pb.PGSHA256Response, error) {
	span := tracing.Tracer.StartSpan("GetPGSHA256Hash")
	defer span.Finish()

	span.SetTag(tagRequestId, getRequestId(ctx))
	return t.wrapped.GetPGSHA256Hash(ctx, req)
}

func (t *requestTracer) GetMYSQLSHA1Hash(ctx context.Context, req *pb.MYSQLSHA1HashRequest) (*pb.MYSQLSHA1Response, error) {
	span := tracing.Tracer.StartSpan("GetMYSQLSHA1Hash")
	defer span.Finish()

	span.SetTag(tagRequestId, getRequestId(ctx))
	return t.wrapped.GetMYSQLSHA1Hash(ctx, req)
}

func jaegerTracer(logger *log.Logger) (opentracing.Tracer, io.Closer, error) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		logger.Warnf("unable to parse environmental jaeger config due to %s, falling back to default config", err)
		cfg = &jaegercfg.Configuration{
			ServiceName: "approzium_authentication_server",
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans: true,
			},
		}
	}
	return cfg.NewTracer(
		jaegercfg.Logger(&jaegerLogger{applicationLogger: logger}),
		// TODO should we add a metrics factory here?
	)
}

// jaegerLogger meets jaeger's Logger interface by fulfilling it with our
// application-level logger. This allows us to format Jaeger's logger to
// match our own.
type jaegerLogger struct {
	applicationLogger *log.Logger
}

func (l *jaegerLogger) Error(msg string) {
	l.applicationLogger.Error(msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	l.applicationLogger.Infof(msg, args)
}
