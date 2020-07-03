package identity

import (
	"fmt"
	"time"

	"github.com/approzium/approzium/authenticator/server/metrics"
	pb "github.com/approzium/approzium/authenticator/server/protos"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/metric"
	"go.opencensus.io/metric/metricdata"
)

type Proof struct {
	AuthType   pb.AuthType
	ClientLang pb.ClientLanguage
	AwsAuth    *pb.AWSAuth
}

type Verified struct {
	AuthType pb.AuthType
	IamArn   string
}

type Verifier interface {
	Get(reqLogger *log.Entry, proof *Proof) (*Verified, error)
	Matches(reqLogger *log.Entry, claimedIdentity string, verifiedIdentity *Verified) (bool, error)
}

func NewVerifier() (Verifier, error) {
	numVerAttempts, err := metrics.Registry.AddInt64Cumulative(
		"total_identity_verification_attempts",
		metric.WithDescription("Total attempts to verify caller identity"),
	)
	if err != nil {
		return nil, err
	}
	numVerAttemptsEntry, err := numVerAttempts.GetEntry()
	if err != nil {
		return nil, err
	}

	numVerFailures, err := metrics.Registry.AddInt64Cumulative(
		"total_identity_verification_failures",
		metric.WithDescription("Total failures to verify caller identity"),
	)
	if err != nil {
		return nil, err
	}
	numVerFailuresEntry, err := numVerFailures.GetEntry()
	if err != nil {
		return nil, err
	}

	numMatchAttempts, err := metrics.Registry.AddInt64Cumulative(
		"total_identity_matching_attempts",
		metric.WithDescription("Total checks of whether the identity a caller claims matches their actual identity"),
	)
	if err != nil {
		return nil, err
	}
	numMatchAttemptsEntry, err := numMatchAttempts.GetEntry()
	if err != nil {
		return nil, err
	}

	numMatchFailures, err := metrics.Registry.AddInt64Cumulative(
		"total_identity_matching_failures",
		metric.WithDescription("Total calls where caller claimed an identity that was not their actual identity"),
	)
	if err != nil {
		return nil, err
	}
	numMatchFailuresEntry, err := numMatchFailures.GetEntry()
	if err != nil {
		return nil, err
	}

	awsReqMilliseconds, err := metrics.Registry.AddInt64Gauge(
		"total_aws_request_milliseconds",
		metric.WithDescription("Total AWS request milliseconds"),
		metric.WithUnit(metricdata.UnitMilliseconds),
	)
	if err != nil {
		return nil, err
	}
	awsReqMillisecondsEntry, err := awsReqMilliseconds.GetEntry()
	if err != nil {
		return nil, err
	}
	return &tracker{
		aws:                &aws{},
		numVerAttempts:     numVerAttemptsEntry,
		numVerFailures:     numVerFailuresEntry,
		numMatchAttempts:   numMatchAttemptsEntry,
		numMatchFailures:   numMatchFailuresEntry,
		awsReqMilliseconds: awsReqMillisecondsEntry,
	}, nil
}

type tracker struct {
	aws Verifier

	numVerAttempts     *metric.Int64CumulativeEntry
	numVerFailures     *metric.Int64CumulativeEntry
	numMatchAttempts   *metric.Int64CumulativeEntry
	numMatchFailures   *metric.Int64CumulativeEntry
	awsReqMilliseconds *metric.Int64GaugeEntry
}

func (t *tracker) Get(reqLogger *log.Entry, proof *Proof) (*Verified, error) {
	t.numVerAttempts.Inc(1)

	var verifiedIdentity *Verified
	var err error
	switch proof.AuthType {
	case pb.AuthType_AWS:
		verifiedIdentity, err = t.aws.Get(reqLogger, proof)
	default:
		return nil, fmt.Errorf("unexpected auth type %d received", proof.AuthType)
	}
	if err != nil {
		t.numVerFailures.Inc(1)
		reqLogger.Warn(fmt.Sprintf("couldn't verify %s: %s", proof, err))
	} else {
		reqLogger.Info(fmt.Sprintf("verified %s", verifiedIdentity))
	}
	return verifiedIdentity, err
}

func (t *tracker) Matches(reqLogger *log.Entry, claimedIdentity string, verifiedIdentity *Verified) (bool, error) {
	t.numMatchAttempts.Inc(1)

	var matches bool
	var err error
	switch verifiedIdentity.AuthType {
	case pb.AuthType_AWS:
		start := time.Now().UTC()
		matches, err = t.aws.Matches(reqLogger, claimedIdentity, verifiedIdentity)
		t.awsReqMilliseconds.Set(time.Now().UTC().Sub(start).Milliseconds())
	default:
		return false, fmt.Errorf("unexpected auth type %d received", verifiedIdentity.AuthType)
	}
	if !matches || err != nil {
		t.numMatchFailures.Inc(1)
	}
	if !matches {
		reqLogger.Warnf("claimed identity of %s doesn't match verified identity of %+v", claimedIdentity, verifiedIdentity)
	}
	if err != nil {
		reqLogger.Warnf("unable to match claimed identity of %s with verified identity of %+v due to %s", claimedIdentity, verifiedIdentity, err)
	}
	return matches, err
}
