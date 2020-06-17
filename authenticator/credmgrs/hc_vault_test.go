package credmgrs

import (
	"os"
	"testing"
)

// For testing, run Vault locally via:
// "$ vault server -dev -dev-root-token-id=root"
//
// Then, in this application's environment, set:
// 		VAULT_ADDR=http://localhost:8200
// 		VAULT_TOKEN=root
//
// Enable Vault KV Version 1:
// "$ vault secrets enable -path=approzium -version=1 kv"
//
// Create a test secret body:
// "$ nano test-secret.json"
// {
//		"password": "asdfghjkl",
//		"iam_roles": [
//			"arn:aws:iam::accountid:role/rolename1",
//			"arn:aws:iam::accountid:role/rolename2"
//		]
//	}
// "$ vault kv put approzium/postgresql://my-username:my-password@localhost:5432 dbuser1=@test-secret.json"
func TestHcVaultCredMgr_Password(t *testing.T) {
	// Ensure we have the necessary test environment.
	if addr := os.Getenv("VAULT_ADDR"); addr == "" {
		t.Skip("skipping because VAULT_ADDR is unset")
	}
	if token := os.Getenv("VAULT_TOKEN"); token == "" {
		t.Skip("skipping because VAULT_TOKEN is unset")
	}

	credMgr, err := newHashiCorpVaultCreds()
	if err != nil {
		t.Fatal(err)
	}
	identity := DBKey{
		IAMArn: "arn:aws:iam::accountid:role/rolename2",
		DBHost: "postgresql://my-username:my-password@localhost",
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
