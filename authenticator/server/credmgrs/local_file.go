package credmgrs

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const secretsFileLocation = "/authenticator/server/testing/secrets.yaml"

// newLocalFileCreds is for dev purposes: read credentials from a local file.
func newLocalFileCreds(logger *log.Logger) (CredentialManager, error) {
	creds := make(map[DBKey]string)
	type secrets []struct {
		Dbhost   string `yaml:"dbhost"`
		Dbport   string `yaml:"dbport"`
		Dbuser   string `yaml:"dbuser"`
		Password string `yaml:"password"`
		IamArn   string `yaml:"iam_arn"`
	}

	// To make sure we can find the secrets file, get the absolute path to
	// the file that called this method. This will always be something like
	// /Users/yourname/go/src/github.com/approzium/approzium/authenticator/server/credmgrs/credmgr.go
	_, filename, _, _ := runtime.Caller(1)

	homeDirPath := strings.TrimSuffix(filename, "/authenticator/server/credmgrs/credmgr.go")
	pathToSecrets := homeDirPath + secretsFileLocation
	logger.Infof("loading secrets at %q", pathToSecrets)

	var devCreds secrets
	yamlFile, err := ioutil.ReadFile(pathToSecrets)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, &devCreds); err != nil {
		return nil, err
	}
	for _, cred := range devCreds {
		key := DBKey{
			IAMArn: replaceEnvVars(cred.IamArn),
			DBHost: replaceEnvVars(cred.Dbhost),
			DBPort: replaceEnvVars(cred.Dbport),
			DBUser: replaceEnvVars(cred.Dbuser),
		}
		creds[key] = cred.Password
		logger.Debugf("added dev credential for host %s", cred.Dbhost)
	}
	logger.Info("secrets loaded, please restart authenticator to load edits")
	return &localFileCredMgr{creds: creds}, nil
}

type localFileCredMgr struct {
	creds map[DBKey]string
}

func (l *localFileCredMgr) Name() string {
	return "local file (dev only)"
}

func (l *localFileCredMgr) Password(_ *log.Entry, identity DBKey) (string, error) {
	creds, ok := l.creds[identity]
	if !ok {
		return "", ErrNotFound
	}
	return creds, nil
}

// replaceEnvVars takes fields that are formatted like ${FOO}, strips
// the $ and brackets, and replaces the env var with its environmental
// value
func replaceEnvVars(field string) string {
	if !strings.HasPrefix(field, "$") {
		// It's not an env var.
		return field
	}
	field = strings.ReplaceAll(field, "${", "")
	field = strings.ReplaceAll(field, "}", "")
	return os.Getenv(field)
}
