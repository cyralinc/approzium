package credmgrs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cyralinc/approzium/authenticator/server/config"
	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

// In Vault, a mount path is the path where a secrets engine has been
// mounted. This code supports mounts that have been added the following way:
// "$ vault secrets enable -path=approzium -version=1 kv"
// Someday we may wish to make this path configurable.
const mountPath = "approzium"

func newHashiCorpVaultCreds(_ *log.Logger, config config.Config) (CredentialManager, error) {
	if config.VaultAddr == "" {
		return &hcVaultCredMgr{}, errors.New("no vault address detected")
	}
	credMgr := &hcVaultCredMgr{
		token:     config.VaultToken,
		tokenPath: config.VaultTokenPath,
		addr:      config.VaultAddr,
	}

	// Check that we're able to communicate with Vault by doing a test read.
	client, err := credMgr.vaultClient()
	if err != nil {
		return &hcVaultCredMgr{}, err
	}
	if _, err := client.Logical().Read(mountPath); err != nil {
		return &hcVaultCredMgr{}, err
	}
	return credMgr, nil
}

type hcVaultCredMgr struct {
	token     string
	tokenPath string
	addr      string
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

	password, err := passwordFromSecret(secret.Data, identity)
	if err != nil {
		return "", err
	}
	return password, nil
}

// vaultClient retrieves a client using either the environmental VAULT_TOKEN,
// or reading the latest token from the token file sink.
func (h *hcVaultCredMgr) vaultClient() (*vault.Client, error) {
	// Only use the token sink if a vault token is not provided
	if h.tokenPath != "" && h.token == "" {
		tokenBytes, err := ioutil.ReadFile(h.tokenPath)
		if err != nil {
			return nil, err
		}
		h.token = string(tokenBytes)
	}

	// This uses a default configuration for Vault. This includes reading
	// Vault's environment variables and setting them as a configuration.
	config := vault.DefaultConfig()
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetAddress(h.addr)
	client.SetToken(h.token)
	return client, err
}
