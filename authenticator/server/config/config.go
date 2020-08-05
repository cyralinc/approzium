package config

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const envVarPrefix = "APPROZIUM_"

type section string

const (
	// For fields that shouldn't be included in the file conf
	sectionExclude section = "exclude"

	// For fields configuring how Approzium listens.
	sectionListener = "listener"

	// For fields configuring logging.
	sectionLogging = "logging"

	// For fields configuring TLS.
	sectionTLS = "tls"

	// For fields configuring a secrets manager.
	sectionSecretsMgr = "secrets_manager"
)

var (
	flagConf = &Config{}

	// Add new fields to the config by simply adding them to the fieldRegistry. This automatically adds support
	// for them in env vars, config files, command-line flags, and defaults.
	fieldRegistry = map[section][]field{
		sectionListener: {
			{name: "host", defaultVal: "127.0.0.1", flagConfField: &flagConf.Host, goFieldName: "Host"},
			{name: "http port", defaultVal: 6000, flagConfField: &flagConf.HTTPPort, goFieldName: "HTTPPort"},
			{name: "grpc port", defaultVal: 6001, flagConfField: &flagConf.GRPCPort, goFieldName: "GRPCPort"},
		},
		sectionLogging: {
			{name: "log level", defaultVal: "info", flagConfField: &flagConf.LogLevel, goFieldName: "LogLevel"},
			{name: "log format", defaultVal: "text", flagConfField: &flagConf.LogFormat, goFieldName: "LogFormat"},
			{name: "log raw", defaultVal: false, flagConfField: &flagConf.LogRaw, goFieldName: "LogRaw"},
		},
		sectionTLS: {
			{name: "disable tls", defaultVal: false, flagConfField: &flagConf.DisableTLS, goFieldName: "DisableTLS"},
			{name: "tls cert path", defaultVal: "", flagConfField: &flagConf.PathToTLSCert, goFieldName: "PathToTLSCert"},
			{name: "tls key path", defaultVal: "", flagConfField: &flagConf.PathToTLSKey, goFieldName: "PathToTLSKey"},
		},
		sectionSecretsMgr: {
			{name: "secrets manager", defaultVal: "", flagConfField: &flagConf.SecretsManager, goFieldName: "SecretsManager"},
			{name: "vault token", defaultVal: "", flagConfField: &flagConf.VaultToken, goFieldName: "VaultToken"},
			{name: "vault token path", defaultVal: "", flagConfField: &flagConf.VaultTokenPath, goFieldName: "VaultTokenPath"},
			{name: "vault addr", defaultVal: "", flagConfField: &flagConf.VaultAddr, goFieldName: "VaultAddr"},
			{name: "aws region", defaultVal: "", flagConfField: &flagConf.AwsRegion, goFieldName: "AwsRegion"},
		},
		sectionExclude: {
			{name: "config", defaultVal: "", flagConfField: &flagConf.ConfigFilePath, goFieldName: "ConfigFilePath"},
			{name: "dev", defaultVal: false, flagConfField: &flagConf.DevMode, goFieldName: "DevMode"},
			{name: "version", defaultVal: false, flagConfField: &flagConf.Version, goFieldName: "Version"},
		},
	}

	allRegisteredFields = func() []field {
		var allFields []field
		for _, sectionFields := range fieldRegistry {
			for _, registeredField := range sectionFields {
				allFields = append(allFields, registeredField)
			}
		}
		return allFields
	}()
)

type field struct {
	// A user-friendly name for the field. Multi-word parameters should be separated by spaces: "fizz buzz".
	name string

	// The default value for the field.
	defaultVal interface{}

	// The field on the flag conf to populate with this value.
	flagConfField interface{}

	// The name of the Go field this should be bound with.
	goFieldName string
}

func init() {
	for _, registeredField := range allRegisteredFields {
		switch fieldOnConf := registeredField.flagConfField.(type) {
		case *bool:
			flag.BoolVar(fieldOnConf, flagName(registeredField.name), registeredField.defaultVal.(bool), "")
		case *int:
			flag.IntVar(fieldOnConf, flagName(registeredField.name), registeredField.defaultVal.(int), "")
		case *string:
			flag.StringVar(fieldOnConf, flagName(registeredField.name), registeredField.defaultVal.(string), "")
		}
	}
}

// Config is an object for storing configuration variables set through
// Approzium's environment.
//
// Please see https://approzium.org/configuration for elaboration upon
// each parameter.
type Config struct {
	Host     string
	HTTPPort int
	GRPCPort int

	DisableTLS    bool
	PathToTLSCert string
	PathToTLSKey  string

	LogLevel  string
	LogFormat string
	LogRaw    bool

	SecretsManager string

	// Vault related
	VaultToken     string
	VaultTokenPath string
	VaultAddr      string

	// AWS secrets manager related
	AwsRegion string

	// Special flags
	ConfigFilePath string
	DevMode        bool
	Version        bool
}

func verify(config *Config) error {
	if !config.DisableTLS {
		if config.PathToTLSCert == "" {
			return errors.New("tls is enabled but no tls cert has been provided")
		}
		if config.PathToTLSKey == "" {
			return errors.New("tls is enabled but no tls key has been provided")
		}
	}
	return nil
}

// Parse returns the parsed config. A pointer is not returned
// because after first parse, the config is immutable.
func Parse() (Config, error) {
	config, err := parseWithPrecedence()
	if err != nil {
		return Config{}, err
	}
	if config.Version {
		// Nothing further to do here.
		return *config, nil
	}
	if config.DevMode {
		config, err = devModeConfig()
		if err != nil {
			return Config{}, err
		}
	}
	if err := verify(config); err != nil {
		return Config{}, err
	}
	// We return a non-pointer to signal that it will never change.
	return *config, nil
}

func parseWithPrecedence() (*Config, error) {
	flag.Parse()

	// Our order of precedence is:
	// 1. Command-line flag
	// 2. Env var
	// 3. Config file
	// 4. Default
	conf, err := defaultConfig()
	if err != nil {
		return nil, err
	}
	if flagConf.ConfigFilePath != "" {
		if err := overrideWithFileConf(conf); err != nil {
			return nil, err
		}
	}
	if err := overrideWithEnvConf(conf); err != nil {
		return nil, err
	}
	if err := overrideWithFlagConf(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func defaultConfig() (*Config, error) {
	conf := &Config{}
	writableConfFields := reflect.ValueOf(conf).Elem()
	for _, registeredField := range allRegisteredFields {
		strDefault := fmt.Sprintf("%v", registeredField.defaultVal)
		if err := setField(writableConfFields, registeredField, strDefault); err != nil {
			return nil, err
		}
	}
	return conf, nil
}

func overrideWithFileConf(conf *Config) error {
	writableConfFields := reflect.ValueOf(conf).Elem()

	body, err := ioutil.ReadFile(flagConf.ConfigFilePath)
	if err != nil {
		return err
	}

	bodyMap := make(map[section]map[string]interface{})
	if err := yaml.Unmarshal(body, &bodyMap); err != nil {
		return err
	}

	for section, registeredFields := range fieldRegistry {

		sectionFields, exists := bodyMap[section]
		if !exists {
			continue
		}

		for _, registeredField := range registeredFields {
			fieldIfc, exists := sectionFields[fileName(registeredField.name)]
			if !exists {
				continue
			}
			strVal := fmt.Sprintf("%v", fieldIfc)
			if err := setField(writableConfFields, registeredField, strVal); err != nil {
				return err
			}
		}
	}
	return nil
}

// overrideWithFlagConf returns any command-line flags that were set in,
// or default values.
func overrideWithFlagConf(conf *Config) error {
	writableConfFields := reflect.ValueOf(conf).Elem()
	for _, registeredField := range allRegisteredFields {
		f, found := getFlag(flagName(registeredField.name))
		if !found {
			continue
		}
		if err := setField(writableConfFields, registeredField, f.Value.String()); err != nil {
			return err
		}
	}
	return nil
}

func overrideWithEnvConf(conf *Config) error {
	writableConfFields := reflect.ValueOf(conf).Elem()
	for _, registeredField := range allRegisteredFields {
		fieldValue := os.Getenv(envVarName(registeredField.name))
		if fieldValue == "" {
			continue
		}
		if err := setField(writableConfFields, registeredField, fieldValue); err != nil {
			return err
		}
	}
	return nil
}

func devModeConfig() (*Config, error) {
	// dev mode uses the local file back-end
	if err := os.Unsetenv("VAULT_ADDR"); err != nil {
		return &Config{}, err
	}
	return &Config{
		Host:       "127.0.0.1",
		HTTPPort:   6000,
		GRPCPort:   6001,
		DisableTLS: true,
		LogLevel:   "debug",
		LogFormat:  "text",
		LogRaw:     false,
		DevMode:    true,
	}, nil
}

func envVarName(name string) string {
	return envVarPrefix + strings.ToUpper(strings.ReplaceAll(name, " ", "_"))
}

func fileName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}

func flagName(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}

func getFlag(name string) (*flag.Flag, bool) {
	toReturn := &flag.Flag{}
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			toReturn = f
			found = true
		}
	})
	return toReturn, found
}

func setField(elem reflect.Value, registeredField field, value string) error {
	field := elem.FieldByName(registeredField.goFieldName)

	if reflect.DeepEqual(field, reflect.Value{}) {
		return fmt.Errorf("unable to locate %s", registeredField.goFieldName)
	}

	switch registeredField.flagConfField.(type) {
	case *bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case *int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))
	case *string:
		field.SetString(value)
	default:
		return fmt.Errorf("unrecognized field type: %T", registeredField.flagConfField)
	}
	return nil
}
