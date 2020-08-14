package credmgrs

import (
	"errors"
	"fmt"

	"github.com/cyralinc/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotAuthorized = errors.New("not authorized")
	ErrNotFound      = errors.New("not found")
)

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

	// List retrieves all the DBKeys in the database, unsorted and in an
	// inconsistent order. Sorting, if needed, is left to the caller.
	// List is currently only used by the CLI, so has no security-oriented
	// logging or metrics. If that changes, they need to be implemented
	// in the tracker.
	List() ([]*DBCred, error)

	// Creates or clobbers any matching dbCred for the given db host, port, and user.
	// Write is currently only used by the CLI, so has no security-oriented
	// logging or metrics. If that changes, they need to be implemented
	// in the tracker.
	Write(toWrite *DBCred) error

	// Delete deletes a given dbCred.
	// Delete is currently only used by the CLI, so has no security-oriented
	// logging or metrics. If that changes, they need to be implemented
	// in the tracker.
	Delete(toDelete *DBCred) error
}

type instantiator func(*log.Logger, config.Config) (CredentialManager, error)

var (
	credMgrs = map[string]instantiator{
		"vault": newHashiCorpVaultCreds,
		"asm":   newAWSSecretManagerCreds,
		"local": newLocalFileCreds,
	}

	credMgrOptions = func() []string {
		var opts []string
		for userKey := range credMgrs {
			opts = append(opts, userKey)
		}
		return opts
	}()
)

func selectCredMgr(logger *log.Logger, c config.Config) (CredentialManager, error) {
	if c.SecretsManager == "" {
		return legacySelectCredMgr(logger, c)
	}
	credMgrNew, ok := credMgrs[c.SecretsManager]
	if !ok {
		msg := fmt.Sprintf("Unknown secrets manager option: %s. Valid options are %+q", c.SecretsManager, credMgrOptions)
		return nil, fmt.Errorf(msg)
	}
	credMgr, err := credMgrNew(logger, c)
	if err != nil {
		msg := fmt.Sprintf("could not configure %s as credential manager due to err: %s", c.SecretsManager, err)
		return nil, errors.New(msg)
	}
	logger.Infof("using %s as credentials manager", credMgr.Name())
	return credMgr, nil
}

func legacySelectCredMgr(logger *log.Logger, c config.Config) (CredentialManager, error) {
	// Legacy behaviour: try vault then local file
	credMgr, err := newHashiCorpVaultCreds(logger, c)
	if err != nil {
		logger.Debugf("didn't select HashiCorp Vault as credential manager due to err: %s", err)
	} else {
		logger.Info("selected HashiCorp Vault as credential manager")
		return credMgr, nil
	}

	credMgr, err = newLocalFileCreds(logger, c)
	if err != nil {
		logger.Debugf("didn't select local file as credential manager due to err: %s", err)
	} else {
		logger.Info("selected local file as credential manager")
		logger.Warn("local file credential manager should not be used in production")
		return credMgr, err
	}
	return nil, errors.New("no valid credential manager available, see debug-level logs for more information")
}
