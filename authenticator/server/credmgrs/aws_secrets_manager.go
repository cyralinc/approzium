package credmgrs

import (
	"fmt"
    "encoding/json"

    "github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)


func newAWSSecretManagerCreds() (CredentialManager, error) {
    sess, err := session.NewSession()
    if err != nil {
        return nil, err
    }
    //Create a Secrets Manager client
	svc := secretsmanager.New(sess, aws.NewConfig())
    if svc == nil {
        return nil, fmt.Errorf("Cannot instantiate AWS Secrets Manager Client")
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

func (a *awsSecretsManagerCredMgr) Password(_ *log.Entry, identity DBKey) (string, error) {
	path := mountPath + "/" + identity.DBHost + "@" + identity.DBPort
    input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(path),
	}
    log.Infof("path=%s", path)
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
