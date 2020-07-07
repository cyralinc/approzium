package server

import (
	"context"
	"time"

	"github.com/approzium/approzium/authenticator/server/metrics"
	pb "github.com/approzium/approzium/authenticator/server/protos"
	"go.opencensus.io/metric"
	"go.opencensus.io/metric/metricdata"
)

func newRequestMetrics(wrapped pb.AuthenticatorServer) (pb.AuthenticatorServer, error) {
	numRequests, err := metrics.Registry.AddInt64Cumulative("total_requests", metric.WithDescription("Total number of requests"))
	if err != nil {
		return nil, err
	}
	numRequestsEntry, err := numRequests.GetEntry()
	if err != nil {
		return nil, err
	}

	numResponses, err := metrics.Registry.AddInt64Cumulative(
		"total_responses",
		metric.WithDescription("Total number of responses"),
	)
	if err != nil {
		return nil, err
	}
	numResponsesEntry, err := numResponses.GetEntry()
	if err != nil {
		return nil, err
	}

	numErrResponses, err := metrics.Registry.AddInt64Cumulative(
		"total_error_responses",
		metric.WithDescription("Total number of error responses"),
	)
	if err != nil {
		return nil, err
	}
	numErrResponsesEntry, err := numErrResponses.GetEntry()
	if err != nil {
		return nil, err
	}

	totalReqMs, err := metrics.Registry.AddInt64Gauge(
		"total_request_milliseconds",
		metric.WithDescription("Total request milliseconds"),
		metric.WithUnit(metricdata.UnitMilliseconds),
	)
	if err != nil {
		return nil, err
	}
	totalReqMsEntry, err := totalReqMs.GetEntry()
	if err != nil {
		return nil, err
	}
	return &requestMetrics{
		wrapped:         wrapped,
		numRequests:     numRequestsEntry,
		numResponses:    numResponsesEntry,
		numErrResponses: numErrResponsesEntry,
		reqMilliseconds: totalReqMsEntry,
	}, nil
}

type requestMetrics struct {
	wrapped pb.AuthenticatorServer

	numRequests     *metric.Int64CumulativeEntry
	numResponses    *metric.Int64CumulativeEntry
	numErrResponses *metric.Int64CumulativeEntry
	reqMilliseconds *metric.Int64GaugeEntry
}

func (m *requestMetrics) GetPGMD5Hash(ctx context.Context, req *pb.PGMD5HashRequest) (*pb.PGMD5Response, error) {
	start := time.Now().UTC()
	m.numRequests.Inc(1)

	resp, err := m.wrapped.GetPGMD5Hash(ctx, req)
	if err != nil {
		m.numErrResponses.Inc(1)
	}

	m.numResponses.Inc(1)
	m.reqMilliseconds.Set(time.Now().UTC().Sub(start).Milliseconds())
	return resp, err
}

func (m *requestMetrics) GetPGSHA256Hash(ctx context.Context, req *pb.PGSHA256HashRequest) (*pb.PGSHA256Response, error) {
	start := time.Now().UTC()
	m.numRequests.Inc(1)

	resp, err := m.wrapped.GetPGSHA256Hash(ctx, req)
	if err != nil {
		m.numErrResponses.Inc(1)
	}

	m.numResponses.Inc(1)
	m.reqMilliseconds.Set(time.Now().UTC().Sub(start).Milliseconds())
	return resp, err
}

func (m *requestMetrics) GetMYSQLSHA1Hash(ctx context.Context, req *pb.MYSQLSHA1HashRequest) (*pb.MYSQLSHA1Response, error) {
	start := time.Now().UTC()
	m.numRequests.Inc(1)

	resp, err := m.wrapped.GetMYSQLSHA1Hash(ctx, req)
	if err != nil {
		m.numErrResponses.Inc(1)
	}

	m.numResponses.Inc(1)
	m.reqMilliseconds.Set(time.Now().UTC().Sub(start).Milliseconds())
	return resp, err
}
