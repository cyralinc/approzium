package credmgrs

import (
	"errors"
	"fmt"

	"github.com/cyralinc/approzium/authenticator/server/config"
	"github.com/cyralinc/approzium/authenticator/server/metrics"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/metric"
	"go.opencensus.io/metric/metricdata"
)

var (
	ErrNotAuthorized = errors.New("not authorized")
	ErrNotFound      = errors.New("not found")
    options = []string{"vault,", "aws", "local"}
)

type DBKey struct {
	IAMArn string
	DBHost string
	DBPort string
	DBUser string
}

// RetrieveConfigured checks the environment for configured cred
// providers, and selects the first working configuration.
func RetrieveConfigured(logger *log.Logger, config config.Config) (CredentialManager, error) {
	credMgr, err := selectCredMgr(logger, config)
	if err != nil {
		return nil, err
	}
	return newTracker(credMgr)
}

type CredentialManager interface {
	// Name should provide a loggable and error-returnable identifying
	// name for the credential manager.
	Name() string

	// Password should retrieve the password for a given identity.
	// If the identity is not found, an error should be returned.
	// IMPORTANT: While the identity given for the password should
	// be trusted, we should not assume the identity should have
	// access to the database they're requesting it for. It's the
	// responsibility of the Password call to ensure that the given
	// IAM ARN _should_ have access to the given DB.
	Password(reqLogger *log.Entry, identity DBKey) (string, error)
}

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

func selectCredMgr(logger *log.Logger, config_ config.Config) (CredentialManager, error) {
	if config_.SecretsManager == "" {
		return legacySelectCredMgr(logger, config_)
	}
	credMgrs := map[string]func(*log.Logger, config.Config) (CredentialManager, error){
		"vault": newHashiCorpVaultCreds,
		"aws":   newAWSSecretManagerCreds,
		"local": newLocalFileCreds,
	}
	credMgrNew, ok := credMgrs[config_.SecretsManager]
	if !ok {
		return nil, fmt.Errorf("Unknown secrets manager option: %s", config_.SecretsManager)
	}
	credMgr, err := credMgrNew(logger, config_)
	if err != nil {
		logger.Debugf("didn't select %s as credential manager due to err: %s", credMgr.Name(), err)
		return nil, errors.New("no valid credential manager available, see debug-level logs for more information")
	}
	logger.Infof("using %s as credentials manager", credMgr.Name())
	return credMgr, nil
}

func legacySelectCredMgr(logger *log.Logger, config_ config.Config) (CredentialManager, error) {
	// Legacy behaviour: try vault then local file
	credMgr, err := newHashiCorpVaultCreds(logger, config_)
	if err != nil {
		logger.Debugf("didn't select HashiCorp Vault as credential manager due to err: %s", err)
	} else {
		logger.Info("selected HashiCorp Vault as credential manager")
		return credMgr, nil
	}

	credMgr, err = newLocalFileCreds(logger, config_)
	if err != nil {
		logger.Debugf("didn't select local file as credential manager due to err: %s", err)
	} else {
		logger.Info("selected local file as credential manager")
		logger.Warn("local file credential manager should not be used in production")
		return credMgr, err
	}
	return nil, errors.New("no valid credential manager available, see debug-level logs for more information")
}

func passwordFromSecret(secret map[string]interface{}, identity DBKey) (string, error) {
	// Please see tests for examples of the kind of secret data we'd expect
	// here.
	userData := secret[identity.DBUser]
	userDataMap, ok := userData.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("couldn't convert %s to a string, type is %T", userData, userData)
	}

	// Verify that the inbound IAM role is one of the IAM roles listed as appropriate.
	iamArnsRaw, ok := userDataMap["iam_arns"]
	if !ok {
		return "", fmt.Errorf("iam_arns not found in %s", userDataMap)
	}
	iamArns, ok := iamArnsRaw.([]interface{})
	if !ok {
		return "", fmt.Errorf("could not convert %s to array, type is %T", iamArnsRaw, iamArnsRaw)
	}
	authorized := false
	for _, iamArnRaw := range iamArns {
		iamArn, ok := iamArnRaw.(string)
		if !ok {
			return "", fmt.Errorf("couldn't convert %s to a string, type is %T", iamArnRaw, iamArnRaw)
		}
		if iamArn == identity.IAMArn {
			authorized = true
			break
		}
	}
	if !authorized {
		return "", ErrNotAuthorized
	}

	// Verification passed. OK to return the password.
	passwordRaw, ok := userDataMap["password"]
	if !ok {
		return "", fmt.Errorf("password not found in %s", userDataMap)
	}
	password, ok := passwordRaw.(string)
	if !ok {
		return "", fmt.Errorf("could not convert %s to string, type is %T", passwordRaw, passwordRaw)
	}
	return password, nil
}
