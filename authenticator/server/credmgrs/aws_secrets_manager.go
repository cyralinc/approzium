package credmgrs

import (
	"encoding/json"
	"fmt"
	"strings"

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

func (a *awsSecretsManagerCredMgr) Password(_ *log.Entry, identity DBKey) (string, error) {
	dbCreds, err := a.read(identity.DBHost, identity.DBPort)
	if err != nil {
		return "", err
	}
	password, err := getPasswordIfAuthorized(dbCreds, identity)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (a *awsSecretsManagerCredMgr) List() ([]*DBCred, error) {
	input := &secretsmanager.ListSecretsInput{}
	result, err := a.Client.ListSecrets(input)
	if err != nil {
		return nil, err
	}
	var all []*DBCred
	for _, entry := range result.SecretList {
		if !strings.HasPrefix(*entry.Name, mountPath) {
			continue
		}
		subPath := strings.TrimPrefix(*entry.Name, mountPath+"/")

		hostPort := strings.Split(subPath, "@")
		if len(hostPort) != 2 {
			return nil, fmt.Errorf("unexpected value: %s", *entry.Name)
		}
		host := hostPort[0]
		port := hostPort[1]

		dbCreds, err := a.read(host, port)
		if err != nil {
			return nil, err
		}
		all = append(all, dbCreds...)
	}
	return all, nil
}

func (a *awsSecretsManagerCredMgr) Write(toWrite *DBCred) error {
	path := asmSecretPath(toWrite.Host, toWrite.Port)

	dbCreds, err := a.read(toWrite.Host, toWrite.Port)
	if err != nil {
		if err != ErrNotFound {
			return err
		}

		// Create the secret, since it doesn't yet exist.
		b, err := json.Marshal(toWrite.UserMap())
		if err != nil {
			return err
		}
		if _, err := a.Client.CreateSecret(&secretsmanager.CreateSecretInput{
			Name:         aws.String(path),
			SecretString: aws.String(fmt.Sprintf("%s", b)),
		}); err != nil {
			return err
		}
		return nil
	}

	// Now add or overwrite the user in the existing map of users.
	secret := toSecret(dbCreds)
	secret[toWrite.User] = toWrite.UserDataMap()

	b, err := json.Marshal(secret)
	if err != nil {
		return err
	}
	if _, err := a.Client.UpdateSecret(&secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(path),
		SecretString: aws.String(fmt.Sprintf("%s", b)),
	}); err != nil {
		return err
	}
	return nil
}

func (a *awsSecretsManagerCredMgr) Delete(toDelete *DBCred) error {
	path := asmSecretPath(toDelete.Host, toDelete.Port)

	// Since multiple db creds are housed at a particular host and
	// port, we need to delete this one and save the remaining. Or,
	// if none remain, we need to delete the entire host/port record.
	dbCreds, err := a.read(toDelete.Host, toDelete.Port)
	if err != nil {
		return err
	}
	dbCreds = deleteIfExists(toDelete, dbCreds)

	if len(dbCreds) == 0 {
		if _, err := a.Client.DeleteSecret(&secretsmanager.DeleteSecretInput{
			SecretId: aws.String(path),
		}); err != nil {
			return err
		}
	} else {
		secret := toSecret(dbCreds)
		b, err := json.Marshal(secret)
		if err != nil {
			return err
		}
		if _, err := a.Client.UpdateSecret(&secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(path),
			SecretString: aws.String(fmt.Sprintf("%s", b)),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (a *awsSecretsManagerCredMgr) read(host, port string) ([]*DBCred, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(asmSecretPath(host, port)),
	}
	result, err := a.Client.GetSecretValue(input)
	if err != nil {
		if _, ok := err.(*secretsmanager.ResourceNotFoundException); ok {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if result.SecretString == nil {
		return nil, fmt.Errorf("no secret string returned from AWS Secrets Manager")
	}
	secretString := *result.SecretString
	secret := make(map[string]interface{})
	if err := json.Unmarshal([]byte(secretString), &secret); err != nil {
		return nil, err
	}
	return toDbCreds(host, port, secret)
}

func asmSecretPath(host, port string) string {
	// AWS Secrets Manager does not support ":" in their secret names
	return mountPath + "/" + host + "@" + port
}
