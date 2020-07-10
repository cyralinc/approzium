package credmgrs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

// In Vault, a mount path is the path where a secrets engine has been
// mounted. This code supports mounts that have been added the following way:
// "$ vault secrets enable -path=approzium -version=1 kv"
// Someday we may wish to make this path configurable.
const mountPath = "approzium"

func newHashiCorpVaultCreds(tokenPath string) (CredentialManager, error) {
	if addr := os.Getenv(vault.EnvVaultAddress); addr == "" {
		return nil, errors.New("no vault address detected")
	}
	credMgr := &hcVaultCredMgr{
		tokenPath: tokenPath,
	}

	// Check that we're able to communicate with Vault by doing a test read.
	client, err := credMgr.vaultClient()
	if err != nil {
		return nil, err
	}
	if _, err := client.Logical().Read(mountPath); err != nil {
		return nil, err
	}
	return credMgr, nil
}

type hcVaultCredMgr struct {
	tokenPath string
}

func (h *hcVaultCredMgr) Name() string {
	return "HashiCorp Vault"
}

func (h *hcVaultCredMgr) Password(_ *log.Entry, identity DBKey) (string, error) {
	client, err := h.vaultClient()
	if err != nil {
		return "", err
	}

	path := mountPath + "/" + identity.DBHost + ":" + identity.DBPort
	secret, err := client.Logical().Read(path)
	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("nothing exists at this Vault path")
	}
	if secret.Data == nil {
		return "", fmt.Errorf("no response body data returned from Vault")
	}

	// Please see tests for examples of the kind of secret data we'd expect
	// here.
	userData := secret.Data[identity.DBUser]
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

// vaultClient retrieves a client using either the environmental VAULT_TOKEN,
// or reading the latest token from the token file sink.
func (h *hcVaultCredMgr) vaultClient() (*vault.Client, error) {
	// Only use the token sink if there's not already an environmental
	// VAULT_TOKEN.
	if h.tokenPath != "" && os.Getenv(vault.EnvVaultToken) == "" {
		tokenBytes, err := ioutil.ReadFile(h.tokenPath)
		if err != nil {
			return nil, err
		}
		// There is no way to directly pass in the token, so we
		// must set it in the environment.
		os.Setenv(vault.EnvVaultToken, string(tokenBytes))
		defer os.Unsetenv(vault.EnvVaultToken)
	}

	// This uses a default configuration for Vault. This includes reading
	// Vault's environment variables and setting them as a configuration.
	return vault.NewClient(nil)
}
