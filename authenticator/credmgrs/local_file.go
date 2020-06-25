package credmgrs

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const envVarTestRole = "TEST_IAM_ROLE"

// newLocalFileCreds is for dev purposes: read credentials from a local file.
func newLocalFileCreds() (CredentialManager, error) {
	creds := make(map[DBKey]string)
	type secrets []struct {
		Dbhost   string `yaml:"dbhost"`
		Dbport   string `yaml:"dbport"`
		Dbuser   string `yaml:"dbuser"`
		Password string `yaml:"password"`
		IamArn   string `yaml:"iam_arn"`
	}
	var devCreds secrets
	yamlFile, err := ioutil.ReadFile("testing/secrets.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &devCreds)
	if err != nil {
		return nil, err
	}
	for _, cred := range devCreds {
		iamArn := cred.IamArn
		if strings.HasPrefix(cred.IamArn, "$") {
			// It's an env var, parse it.
			iamArn = strings.ReplaceAll(iamArn, "${", "")
			iamArn = strings.ReplaceAll(iamArn, "}", "")
			iamArn = os.Getenv(iamArn)
		}
		key := DBKey{
			IAMArn: iamArn,
			DBHost: cred.Dbhost,
			DBPort: cred.Dbport,
			DBUser: cred.Dbuser,
		}
		creds[key] = cred.Password
		log.Debugf("added dev credential for host %s", cred.Dbhost)
	}
	return &localFileCredMgr{creds: creds}, nil
}

type localFileCredMgr struct {
	creds map[DBKey]string
}

func (l *localFileCredMgr) Name() string {
	return "local file (dev only)"
}

func (l *localFileCredMgr) Password(identity DBKey) (string, error) {
	creds, ok := l.creds[identity]
	if !ok {
		return "", ErrNotFound
	}
	return creds, nil
}
