package metrics

import (
	"go.opencensus.io/metric"
	"go.opencensus.io/metric/metricproducer"
)

var Registry = func() *metric.Registry {
	r := metric.NewRegistry()
	metricproducer.GlobalManager().AddProducer(r)
	return r
}()
