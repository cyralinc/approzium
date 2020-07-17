package credmgrs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/cyralinc/approzium/authenticator/server/config"
	"testing"
)

func TestAwsSecretsManager(t *testing.T) {
	sess, err := session.NewSession()
	if err != nil {
		t.Skip("skipping because no AWS config is available")
	}

	svc := secretsmanager.New(sess, aws.NewConfig())

	// Sunny path: plant a credential and retrieve it through Approzium
	// for these tests to work, the secret has to already exist on Asm, so the tests populate them with values
	identity := DBKey{
		IAMArn: "arn:aws:iam::accountid:role/rolename2",
		DBHost: "127.0.0.1",
		DBPort: "5432",
		DBUser: "dbuser1",
	}
	path := AsmSecretPath(identity)
	secretString := `
{
    "dbuser1": {
        "password": "asdfghjkl",
        "iam_arns": [
            "arn:aws:iam::accountid:role/rolename1",
            "arn:aws:iam::accountid:role/rolename2"
        ]
    }
}`
	input := &secretsmanager.PutSecretValueInput{
		SecretId:     &path,
		SecretString: &secretString,
	}
	result, err := svc.PutSecretValue(input)
	if err != nil {
		t.Fatal(err)
	}

	credMgr, err := newAWSSecretManagerCreds(nil, config.Config{})
	if err != nil {
		t.Fatal(err)
	}
	password, err := credMgr.Password(testLogEntry, identity)
	if err != nil {
		t.Fatal(err)
	}
	if password != "asdfghjkl" {
		t.Fatalf("expected: %s; actual: %s", "asdfghjkl", password)
	}

	// Remove access and ensure it denies access
	secretString = `
{
    "dbuser1": {
        "password": "asdfghjkl",
        "iam_arns": [
            "arn:aws:iam::accountid:role/rolename1"
        ]
    }
}`
	input = &secretsmanager.PutSecretValueInput{
		SecretId:     &path,
		SecretString: &secretString,
	}
	result, err = svc.PutSecretValue(input)
	if err != nil {
		t.Fatal(err)
	}
	password, err = credMgr.Password(testLogEntry, identity)

	if err != ErrNotAuthorized {
		t.Fatal(err)
	}

	_ = result
}
