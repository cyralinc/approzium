package credmgrs

import (
	"github.com/cyralinc/approzium/authenticator/server/metrics"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/metric"
	"go.opencensus.io/metric/metricdata"
)

func newTracker(wrapped CredentialManager) (CredentialManager, error) {
	numPwAttempts, err := metrics.Registry.AddInt64Cumulative(
		"total_password_retrieval_attempts",
		metric.WithDescription("The number of times a caller has requested a password from the database to authenticate"),
	)
	if err != nil {
		return nil, err
	}
	numPwAttemptsEntry, err := numPwAttempts.GetEntry()
	if err != nil {
		return nil, err
	}

	numPwFailures, err := metrics.Registry.AddInt64Cumulative(
		"total_password_retrieval_failures",
		metric.WithDescription("The number of times a caller has failed to retrieve a password for any reason"),
	)
	if err != nil {
		return nil, err
	}
	numPwFailuresEntry, err := numPwFailures.GetEntry()
	if err != nil {
		return nil, err
	}

	numPwUnauthorized, err := metrics.Registry.AddInt64Cumulative(
		"total_password_retrieval_unauthorized",
		metric.WithDescription("The number of times a caller has been unauthorized to retrieve a password"),
	)
	if err != nil {
		return nil, err
	}
	numPwUnauthorizedEntry, err := numPwUnauthorized.GetEntry()
	if err != nil {
		return nil, err
	}

	pwReqMilliseconds, err := metrics.Registry.AddInt64Gauge(
		"total_password_request_milliseconds",
		metric.WithDescription("Total password retrieval milliseconds"),
		metric.WithUnit(metricdata.UnitMilliseconds),
	)
	if err != nil {
		return nil, err
	}
	pwReqMillisecondsEntry, err := pwReqMilliseconds.GetEntry()
	if err != nil {
		return nil, err
	}
	return &tracker{
		wrapped:           wrapped,
		numPwAttempts:     numPwAttemptsEntry,
		numPwFailures:     numPwFailuresEntry,
		numPwUnauthorized: numPwUnauthorizedEntry,
		pwReqMilliseconds: pwReqMillisecondsEntry,
	}, nil
}

type tracker struct {
	wrapped CredentialManager

	numPwAttempts     *metric.Int64CumulativeEntry
	numPwFailures     *metric.Int64CumulativeEntry
	numPwUnauthorized *metric.Int64CumulativeEntry
	pwReqMilliseconds *metric.Int64GaugeEntry
}

func (t *tracker) Name() string {
	return t.wrapped.Name()
}

func (t *tracker) Password(reqLogger *log.Entry, identity DBKey) (string, error) {
	t.numPwAttempts.Inc(1)

	password, err := t.wrapped.Password(reqLogger, identity)
	if err != nil {
		t.numPwFailures.Inc(1)
		reqLogger.Warnf("failed attempt to retrieve identity %+v due to %s", identity, err)
		if err == ErrNotAuthorized {
			t.numPwUnauthorized.Inc(1)
		}
	}
	return password, err
}

func (t *tracker) List() ([]*DBCred, error) {
	return t.wrapped.List()
}

func (t *tracker) Write(toWrite *DBCred) error {
	return t.wrapped.Write(toWrite)
}

func (t *tracker) Delete(toDelete *DBCred) error {
	return t.wrapped.Delete(toDelete)
}
