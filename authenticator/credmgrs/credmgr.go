package credmgrs

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	ErrNotAuthorized = errors.New("not authorized")
	ErrNotFound      = errors.New("not found")
)

type DBKey struct {
	IAMArn string
	DBHost string
	DBPort string
	DBUser string
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
	Password(identity DBKey) (string, error)
}

// RetrieveConfigured checks the environment for configured cred
// providers, and selects the first working configuration.
func RetrieveConfigured() (CredentialManager, error) {
	credMgr, err := newHashiCorpVaultCreds()
	if err != nil {
		log.Debugf("didn't select HashiCorp Vault as credential manager due to err: %s", err)
	} else {
		log.Info("selected HashiCorp Vault as credential manager")
		return credMgr, nil
	}

	credMgr, err = newLocalFileCreds()
	if err != nil {
		log.Debugf("didn't select local file as credential manager due to err: %s", err)
	} else {
		log.Info("selected local file as credential manager")
		return credMgr, err
	}
	return nil, errors.New("no valid credential manager available, see debug-level logs for more information")
}
