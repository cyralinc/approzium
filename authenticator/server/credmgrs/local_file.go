package credmgrs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cyralinc/approzium/authenticator/server/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const secretsFileLocation = "/authenticator/server/testing/secrets.yaml"

// newLocalFileCreds is for dev purposes: read credentials from a local file.
func newLocalFileCreds(logger *log.Logger, c config.Config) (CredentialManager, error) {
	pathToSecrets := c.LocalFilePath
	if pathToSecrets == "" {
		// To make sure we can find the secrets file, get the absolute path to
		// the file that called this method. This will always be something like
		// /Users/yourname/go/src/github.com/approzium/approzium/authenticator/server/credmgrs/credmgr.go
		_, filename, _, _ := runtime.Caller(1)

		homeDirPath := strings.TrimSuffix(filename, "/authenticator/server/credmgrs/credmgr.go")
		pathToSecrets = homeDirPath + secretsFileLocation
	}
	logger.Infof("loading secrets at %q", pathToSecrets)

	logger.Warn("local file credential manager should not be used in production")
	return &localFileCredMgr{pathToSecrets: pathToSecrets}, nil
}

type localFileCredMgr struct {
	pathToSecrets string
}

func (l *localFileCredMgr) Name() string {
	return "local file (dev only)"
}

func (l *localFileCredMgr) Password(_ *log.Entry, identity DBKey) (string, error) {
	creds, err := l.read()
	if err != nil {
		return "", err
	}
	password, ok := creds[identity]
	if !ok {
		return "", ErrNotFound
	}
	return password, nil
}

func (l *localFileCredMgr) List() ([]*DBCred, error) {
	creds, err := l.read()
	if err != nil {
		return nil, err
	}

	var dbCreds []*DBCred
	for dbKey, password := range creds {
		dbCreds = append(dbCreds, &DBCred{
			Host:            dbKey.DBHost,
			Port:            dbKey.DBPort,
			User:            dbKey.DBUser,
			Password:        password,
			AccessGrantedTo: []string{dbKey.IAMArn},
		})
	}
	return dbCreds, nil
}

func (l *localFileCredMgr) Write(toWrite *DBCred) error {
	creds, err := l.read()
	if err != nil {
		return err
	}

	for _, arn := range toWrite.AccessGrantedTo {
		dbKey := DBKey{
			IAMArn: arn,
			DBHost: toWrite.Host,
			DBPort: toWrite.Port,
			DBUser: toWrite.User,
		}
		creds[dbKey] = toWrite.Password
	}
	return l.write(creds)
}

func (l *localFileCredMgr) Delete(toDelete *DBCred) error {
	creds, err := l.read()
	if err != nil {
		return err
	}

	for _, arnToDelete := range toDelete.AccessGrantedTo {
		dbKey := DBKey{
			IAMArn: arnToDelete,
			DBHost: toDelete.Host,
			DBPort: toDelete.Port,
			DBUser: toDelete.User,
		}
		delete(creds, dbKey)
	}
	return l.write(creds)
}

func (l *localFileCredMgr) read() (map[DBKey]string, error) {
	yamlFile, err := ioutil.ReadFile(filepath.Clean(l.pathToSecrets))
	if err != nil {
		return nil, err
	}

	var devCreds []*localFileSecret
	if err = yaml.Unmarshal(yamlFile, &devCreds); err != nil {
		return nil, err
	}

	creds := make(map[DBKey]string)
	for _, devCred := range devCreds {
		key := DBKey{
			IAMArn: replaceEnvVars(devCred.IamArn),
			DBHost: replaceEnvVars(devCred.Dbhost),
			DBPort: replaceEnvVars(devCred.Dbport),
			DBUser: replaceEnvVars(devCred.Dbuser),
		}
		creds[key] = devCred.Password
	}
	return creds, nil
}

func (l *localFileCredMgr) write(creds map[DBKey]string) error {
	var devCreds []*localFileSecret
	for dbKey, password := range creds {
		devCreds = append(devCreds, &localFileSecret{
			Dbhost:   dbKey.DBHost,
			Dbport:   dbKey.DBPort,
			Dbuser:   dbKey.DBUser,
			Password: password,
			IamArn:   dbKey.IAMArn,
		})
	}

	b, err := yaml.Marshal(devCreds)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Clean(l.pathToSecrets), b, 644)
}

type localFileSecret struct {
	Dbhost   string `yaml:"dbhost"`
	Dbport   string `yaml:"dbport"`
	Dbuser   string `yaml:"dbuser"`
	Password string `yaml:"password"`
	IamArn   string `yaml:"iam_arn"`
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
