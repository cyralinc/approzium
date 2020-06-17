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
	// Password should retrieve the password for a given identity.
	// If the identity is not found, an error should be returned.
	Password(identity DBKey) (string, error)
}

// RetrieveConfigured checks the environment for configured cred
// providers, and selects the first working configuration.
func RetrieveConfigured() (CredentialManager, error) {
	credMgr, err := newHashiCorpVaultCreds()
	if err != nil {
		log.Debugf("didn't select HashiCorp Vault as credential manager due to err: %s", err)
	} else {
		return credMgr, nil
	}

	credMgr, err = newLocalFileCreds()
	if err != nil {
		log.Debugf("didn't select local file as credential manager due to err: %s", err)
	} else {
		return credMgr, err
	}
	return nil, errors.New("no valid credential manager available, see debug-level logs for more information")
}
