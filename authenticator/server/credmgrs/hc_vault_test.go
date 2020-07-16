package credmgrs

import (
	"io/ioutil"
	"os"
	"testing"
	"github.com/cyralinc/approzium/authenticator/server/config"

	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

var testLogEntry = func() *log.Entry {
	logEntry := log.WithFields(log.Fields{"test": "logger"})
	logEntry.Level = log.FatalLevel
	return logEntry
}()

// To run these tests, first start Vault locally like so:
// "$ vault server -dev -dev-root-token-id=root"
//
// Then, in this application's environment, set:
// 		VAULT_ADDR=http://localhost:8200
// 		VAULT_TOKEN=root
func TestHcVaultCredMgr_WithTokenPath(t *testing.T) {
	// Ensure we have the necessary test environment.
	addr := os.Getenv(vault.EnvVaultAddress)
	if addr == "" {
		t.Skip("skipping because VAULT_ADDR is unset")
	}

	// Record the original
	currentToken := os.Getenv(vault.EnvVaultToken)

	tmpFile, err := ioutil.TempFile(os.TempDir(), "approzium-testing-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	text := []byte(currentToken)
	if _, err = tmpFile.Write(text); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Plant the necessary secrets for this test.
	client, err := vault.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	// We'll try mounting the secrets engine we need, but it might already be
	// mounted, so if there's an error mounting it we'll just ignore it.
	_, _ = client.Logical().Write("/sys/mounts/"+mountPath, map[string]interface{}{
		"type": "kv",
		"options": map[string]interface{}{
			"version": "1",
		},
	})
	if _, err := client.Logical().Write("approzium/127.0.0.1:5432", map[string]interface{}{
		"dbuser1": map[string]interface{}{
			"password": "asdfghjkl",
			"iam_arns": []string{
				"arn:aws:iam::accountid:role/rolename1",
				"arn:aws:iam::accountid:role/rolename2",
			},
		},
	}); err != nil {
		t.Fatal(err)
	}

	// Unset the environmental Vault token so it must be pulled through the token sink.
	defer func() {
		os.Setenv(vault.EnvVaultToken, currentToken)
	}()
	os.Unsetenv(vault.EnvVaultToken)

	// Try to read the creds through the cred manager.
	credMgr, err := newHashiCorpVaultCreds(log, config.Config {
        VaultTokenPath: tmpFile.Name(),
    })
	if err != nil {
		t.Fatal(err)
	}
	identity := DBKey{
		IAMArn: "arn:aws:iam::accountid:role/rolename2",
		DBHost: "127.0.0.1",
		DBPort: "5432",
		DBUser: "dbuser1",
	}
	password, err := credMgr.Password(testLogEntry, identity)
	if err != nil {
		t.Fatal(err)
	}
	if password != "asdfghjkl" {
		t.Fatalf("expected: %s; actual: %s", "asdfghjkl", password)
	}
}

func TestHcVaultCredMgr_WithoutTokenPath(t *testing.T) {
	// Ensure we have the necessary test environment.
	addr := os.Getenv(vault.EnvVaultAddress)
	if addr == "" {
		t.Skip("skipping because VAULT_ADDR is unset")
	}

	// Plant the necessary secrets for this test.
	client, err := vault.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	// We'll try mounting the secrets engine we need, but it might already be
	// mounted, so if there's an error mounting it we'll just ignore it.
	_, _ = client.Logical().Write("/sys/mounts/"+mountPath, map[string]interface{}{
		"type": "kv",
		"options": map[string]interface{}{
			"version": "1",
		},
	})
	if _, err := client.Logical().Write("approzium/localhost:5432", map[string]interface{}{
		"dbuser1": map[string]interface{}{
			"password": "asdfghjkl",
			"iam_arns": []string{
				"arn:aws:iam::accountid:role/rolename1",
				"arn:aws:iam::accountid:role/rolename2",
			},
		},
	}); err != nil {
		t.Fatal(err)
	}

	// Try to read the creds through the cred manager.
	credMgr, err := newHashiCorpVaultCreds(log, config.Config {
        VaultTokenPath: "",
    })
	if err != nil {
		t.Fatal(err)
	}
	identity := DBKey{
		IAMArn: "arn:aws:iam::accountid:role/rolename2",
		DBHost: "localhost",
		DBPort: "5432",
		DBUser: "dbuser1",
	}
	password, err := credMgr.Password(testLogEntry, identity)
	if err != nil {
		t.Fatal(err)
	}
	if password != "asdfghjkl" {
		t.Fatalf("expected: %s; actual: %s", "asdfghjkl", password)
	}
}
