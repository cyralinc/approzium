package credmgrs

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/cyralinc/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
)

func newAWSSecretManagerCreds(_ *log.Logger, c config.Config) (CredentialManager, error) {
	if c.AwsRegion == "" {
		return &awsSecretsManagerCredMgr{}, fmt.Errorf("AWS region not set")
	}
	sess, err := session.NewSession()
	if err != nil {
		return &awsSecretsManagerCredMgr{}, err
	}

	// Create an AWS Secrets Manager client
	var svc *secretsmanager.SecretsManager
	if c.AssumeAWSRole != "" {
		creds := stscreds.NewCredentials(sess, c.AssumeAWSRole)
		svc = secretsmanager.New(sess, aws.NewConfig().WithRegion(c.AwsRegion).WithCredentials(creds))
	} else {
		svc = secretsmanager.New(sess, aws.NewConfig().WithRegion(c.AwsRegion))
	}
	if svc == nil {
		return &awsSecretsManagerCredMgr{}, fmt.Errorf("cannot instantiate AWS Secrets Manager Client")
	}

	credMgr := &awsSecretsManagerCredMgr{
		svc,
	}
	return credMgr, nil
}

type awsSecretsManagerCredMgr struct {
	Client *secretsmanager.SecretsManager
}

func (a *awsSecretsManagerCredMgr) Name() string {
	return "AWS Secrets Manager"
}

func AsmSecretPath(identity DBKey) string {
	// AWS Secrets Manager does not support ":" in their secret names
	return mountPath + "/" + identity.DBHost + "@" + identity.DBPort
}

func (a *awsSecretsManagerCredMgr) Password(_ *log.Entry, identity DBKey) (string, error) {
	path := AsmSecretPath(identity)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(path),
	}
	result, err := a.Client.GetSecretValue(input)
	if err != nil {
		return "", err
	}
	if result.SecretString == nil {
		return "", fmt.Errorf("no secret string returned from Vault")
	}
	secretString := *result.SecretString
	secret := make(map[string]interface{})
	err = json.Unmarshal([]byte(secretString), &secret)
	if err != nil {
		return "", err
	}
	password, err := passwordFromSecret(secret, identity)
	if err != nil {
		return "", err
	}
	return password, nil
}
