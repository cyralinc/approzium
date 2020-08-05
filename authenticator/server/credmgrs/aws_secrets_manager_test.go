// IMPORTANT: for these tests to work:
// - the environment variable APPROZIUM_AWS_REGION has to be set. This is necessary for the Go SDK to access AWS Secrets Manager.
package credmgrs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/cyralinc/approzium/authenticator/server/config"
	"os"
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
	// Check if the test secret is there from prior test runs, if it is, modify it. Otherwise, create a new secret
	_, err = svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &path,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				_, err = svc.CreateSecret(&secretsmanager.CreateSecretInput{
					Name:         &path,
					SecretString: &secretString,
				})
				if err != nil {
					t.Fatal(err)
				}
			default:
				t.Fatal(err)
			}
		} else {
			t.Fatal(err)
		}
	}

	input := &secretsmanager.PutSecretValueInput{
		SecretId:     &path,
		SecretString: &secretString,
	}
	result, err := svc.PutSecretValue(input)
	if err != nil {
		t.Fatal(err)
	}

	credMgr, err := newAWSSecretManagerCreds(nil, config.Config{
		AwsRegion: os.Getenv("APPROZIUM_AWS_REGION"),
	})
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
