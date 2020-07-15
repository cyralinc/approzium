package tracing

import "github.com/opentracing/opentracing-go"

var Tracer = opentracing.GlobalTracer()
