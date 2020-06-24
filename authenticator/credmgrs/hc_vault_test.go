package credmgrs

import (
	"os"
	"testing"

	vault "github.com/hashicorp/vault/api"
)

// To run this test, first start Vault locally like so:
// "$ vault server -dev -dev-root-token-id=root"
//
// Then, in this application's environment, set:
// 		VAULT_ADDR=http://localhost:8200
// 		VAULT_TOKEN=root
func TestHcVaultCredMgr_Password(t *testing.T) {
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
	credMgr, err := newHashiCorpVaultCreds()
	if err != nil {
		t.Fatal(err)
	}
	identity := DBKey{
		IAMArn: "arn:aws:iam::accountid:role/rolename2",
		DBHost: "localhost",
		DBPort: "5432",
		DBUser: "dbuser1",
	}
	password, err := credMgr.Password(identity)
	if err != nil {
		t.Fatal(err)
	}
	if password != "asdfghjkl" {
		t.Fatalf("expected: %s; actual: %s", "asdfghjkl", password)
	}
}
