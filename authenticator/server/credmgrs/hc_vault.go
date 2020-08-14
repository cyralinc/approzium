package credmgrs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

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
	dbCreds, err := h.read(identity.DBHost, identity.DBPort)
	if err != nil {
		return "", err
	}
	password, err := getPasswordIfAuthorized(dbCreds, identity)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (h *hcVaultCredMgr) List() ([]*DBCred, error) {
	client, err := h.vaultClient()
	if err != nil {
		return nil, err
	}
	secret, err := client.Logical().List(mountPath)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}
	keysIfc, ok := secret.Data["keys"]
	if !ok {
		return nil, nil
	}
	keys, ok := keysIfc.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to convert %s to list", keysIfc)
	}

	var all []*DBCred
	for _, key := range keys {
		path, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("unable to convert %s to path", key)
		}

		hostPort := strings.Split(path, ":")
		if len(hostPort) != 2 {
			return nil, fmt.Errorf("unexpected value: %s", path)
		}
		host := hostPort[0]
		port := hostPort[1]

		dbCreds, err := h.read(host, port)
		if err != nil {
			return nil, err
		}
		all = append(all, dbCreds...)
	}
	return all, nil
}

func (h *hcVaultCredMgr) Write(toWrite *DBCred) error {
	path := vaultPath(toWrite.Host, toWrite.Port)

	secret := make(map[string]interface{})
	dbCreds, err := h.read(toWrite.Host, toWrite.Port)
	if err != nil {
		if err != ErrNotFound {
			return err
		}
		secret = toWrite.UserMap()
	} else {
		secret = toSecret(dbCreds)
		secret[toWrite.User] = toWrite.UserDataMap()
	}

	client, err := h.vaultClient()
	if err != nil {
		return err
	}
	if _, err := client.Logical().Write(path, secret); err != nil {
		if !isErrNotMounted(err) {
			return err
		}
		// Mount it and try again.
		if err := client.Sys().Mount(mountPath, &vault.MountInput{
			Type:        "kv",
			Description: "Approzium secrets",
		}); err != nil {
			return err
		}
		if _, err := client.Logical().Write(path, secret); err != nil {
			return err
		}
	}
	return nil
}

func (h *hcVaultCredMgr) Delete(toDelete *DBCred) error {
	client, err := h.vaultClient()
	if err != nil {
		return err
	}

	path := vaultPath(toDelete.Host, toDelete.Port)

	// Since multiple db creds are housed at a particular host and
	// port, we need to delete this one and save the remaining. Or,
	// if none remain, we need to delete the entire host/port record.
	dbCreds, err := h.read(toDelete.Host, toDelete.Port)
	if err != nil {
		return err
	}
	dbCreds = deleteIfExists(toDelete, dbCreds)

	if len(dbCreds) == 0 {
		if _, err := client.Logical().Delete(path); err != nil {
			return err
		}
	} else {
		secret := toSecret(dbCreds)
		if _, err := client.Logical().Write(path, secret); err != nil {
			return err
		}
	}
	return nil
}

// vaultClient retrieves a client using either the environmental VAULT_TOKEN,
// or reading the latest token from the token file sink.
func (h *hcVaultCredMgr) vaultClient() (*vault.Client, error) {
	// Only use the token sink if a vault token is not provided
	if h.tokenPath != "" && h.token == "" {
		tokenBytes, err := ioutil.ReadFile(filepath.Clean(h.tokenPath))
		if err != nil {
			return nil, err
		}
		h.token = string(tokenBytes)
	}

	// This uses a default configuration for Vault. This includes reading
	// Vault's environment variables and setting them as a configuration.
	c := vault.DefaultConfig()
	client, err := vault.NewClient(c)
	if err != nil {
		return nil, err
	}
	if err := client.SetAddress(h.addr); err != nil {
		return nil, err
	}
	client.SetToken(h.token)
	return client, err
}

func (h *hcVaultCredMgr) read(host, port string) ([]*DBCred, error) {
	path := vaultPath(host, port)
	client, err := h.vaultClient()
	if err != nil {
		return nil, err
	}
	secret, err := client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, ErrNotFound
	}
	if secret.Data == nil {
		return nil, ErrNotFound
	}
	return toDbCreds(host, port, secret.Data)
}

func vaultPath(host, port string) string {
	return mountPath + "/" + host + ":" + port
}

func isErrNotMounted(err error) bool {
	return strings.Contains(err.Error(), "no handler for route")
}
