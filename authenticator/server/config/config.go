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
	sectionSecretsMgr = "secrets"
)

var (
	flagConf = &Config{}

	// Add new fields to the config by simply adding them to the fieldRegistry. This automatically adds support
	// for them in env vars, config files, command-line flags, and defaults.
	fieldRegistry = map[section][]field{
		sectionListener: {
			{name: "host", defaultVal: "127.0.0.1", flagConfField: &flagConf.Host, goFieldName: "Host", prependEnvVar: true,
				description: "The address at which Approzium's authenticator should run."},
			{name: "http port", defaultVal: 6000, flagConfField: &flagConf.HTTPPort, goFieldName: "HTTPPort", prependEnvVar: true,
				description: "The port that should serve HTTP traffic."},
			{name: "grpc port", defaultVal: 6001, flagConfField: &flagConf.GRPCPort, goFieldName: "GRPCPort", prependEnvVar: true,
				description: "The port that should serve GRPC traffic, used for client authentication requests."},
		},
		sectionLogging: {
			{name: "log level", defaultVal: "info", flagConfField: &flagConf.LogLevel, goFieldName: "LogLevel", prependEnvVar: true,
				description: `Valid values are "panic", "fatal", "error", "warn", "info", "debug", and "trace".`},
			{name: "log format", defaultVal: "text", flagConfField: &flagConf.LogFormat, goFieldName: "LogFormat", prependEnvVar: true,
				description: `Valid values are "text" and "json".`},
			{name: "log raw", defaultVal: false, flagConfField: &flagConf.LogRaw, goFieldName: "LogRaw", prependEnvVar: true,
				description: "Set to true to disable removing sensitive values from logging."},
		},
		sectionTLS: {
			{name: "disable tls", defaultVal: false, flagConfField: &flagConf.DisableTLS, goFieldName: "DisableTLS", prependEnvVar: true,
				description: "Set to true to disable TLS encryption. Not recommended in production environments."},
			{name: "tls cert path", defaultVal: "", flagConfField: &flagConf.PathToTLSCert, goFieldName: "PathToTLSCert", prependEnvVar: true,
				description: "The path to the certificate proving the Approzium authenticator's identity."},
			{name: "tls key path", defaultVal: "", flagConfField: &flagConf.PathToTLSKey, goFieldName: "PathToTLSKey", prependEnvVar: true,
				description: "The path to the private key to be used for decrypting inbound communication."},
		},
		sectionSecretsMgr: {
			{name: "secrets manager", defaultVal: "", flagConfField: &flagConf.SecretsManager, goFieldName: "SecretsManager", prependEnvVar: true,
				description: `Valid values include "vault" for HashiCorp Vault, "asm" for AWS Secrets Manager", or "local" for a file-based development environment.`},
			{name: "vault token", defaultVal: "", flagConfField: &flagConf.VaultToken, goFieldName: "VaultToken", prependEnvVar: false,
				description: `If "vault" is selected as the secrets manager, a Vault token that won't expire for the Approzium authenticator's lifetime. Use "vault token path" instead in production.`},
			{name: "vault token path", defaultVal: "", flagConfField: &flagConf.VaultTokenPath, goFieldName: "VaultTokenPath", prependEnvVar: true,
				description: `If "vault" is selected as the secrets manager, the path to a file containing a Vault token that is constantly refreshed by the Vault agent.`},
			{name: "vault addr", defaultVal: "", flagConfField: &flagConf.VaultAddr, goFieldName: "VaultAddr", prependEnvVar: false,
				description: `If "vault" is selected as the secrets manager, Vault's address.`},
			{name: "assume aws role", defaultVal: "", flagConfField: &flagConf.AssumeAWSRole, goFieldName: "AssumeAWSRole", prependEnvVar: true,
				description: `If "asm" is selected as the secrets manager, an optional role for the Approzium authenticator to assume when communicating with AWS.`},
			{name: "aws region", defaultVal: "", flagConfField: &flagConf.AwsRegion, goFieldName: "AwsRegion", prependEnvVar: false,
				description: `If "asm" is selected as the secrets manager, the region where the AWS Secrets Manager resides.`},
		},
		sectionExclude: {
			{name: "config", defaultVal: "", flagConfField: &flagConf.ConfigFilePath, goFieldName: "ConfigFilePath", prependEnvVar: true,
				description: "The path to a yaml file bearing configuration settings."},
			{name: "dev", defaultVal: false, flagConfField: &flagConf.DevMode, goFieldName: "DevMode", prependEnvVar: true,
				description: `When true, starts the Approzium authenticator in dev mode, with TLS disabled, with the "local" secrets manager, and debug-level logs.`},
			{name: "version", defaultVal: false, flagConfField: &flagConf.Version, goFieldName: "Version", prependEnvVar: true,
				description: `Outputs the Approzium authenticator's version.`},
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

	// A user-friendly description for the field.
	description string

	// The default value for the field.
	defaultVal interface{}

	// The field on the flag conf to populate with this value.
	flagConfField interface{}

	// The name of the Go field this should be bound with.
	goFieldName string

	// Whether to prepend this option's env var with APPROZIUM_ when looking for it.
	prependEnvVar bool
}

func init() {
	for _, registeredField := range allRegisteredFields {
		switch fieldOnConf := registeredField.flagConfField.(type) {
		case *bool:
			flag.BoolVar(fieldOnConf, flagName(registeredField.name), registeredField.defaultVal.(bool), registeredField.description)
		case *int:
			flag.IntVar(fieldOnConf, flagName(registeredField.name), registeredField.defaultVal.(int), registeredField.description)
		case *string:
			flag.StringVar(fieldOnConf, flagName(registeredField.name), registeredField.defaultVal.(string), registeredField.description)
		}
	}

	// Customize command-line help output with our own grouped format.
	flag.Usage = commandLineUsage
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
	AwsRegion     string
	AssumeAWSRole string

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
		if section == sectionExclude {
			continue
		}
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
		fieldValue := os.Getenv(envVarName(registeredField))
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

func envVarName(f field) string {
	envVar := strings.ToUpper(strings.ReplaceAll(f.name, " ", "_"))
	if !f.prependEnvVar {
		return envVar
	}
	return envVarPrefix + envVar
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

// commandLineUsage customizes the response when a user makes a mistake or needs help using flags.
// Normally, flags are outputted as one group, alphabetically by field name. However, we would prefer
// to group them by section and keep them in the same order as they've been added to the registry.
func commandLineUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", "authenticator")
	sectionOrder := []section{sectionListener, sectionTLS, sectionLogging, sectionSecretsMgr, sectionExclude}

	for _, section := range sectionOrder {
		sectionFields := fieldRegistry[section]
		if section != sectionExclude {
			fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("\n%s", section), "\n")
		} else {
			fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("\n%s", "miscellaneous"), "\n")
		}

		for _, field := range sectionFields {
			switch f := field.flagConfField.(type) {
			case *bool:
				fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("  -%s %T", flagName(field.name), *f), "\n")
			case *int:
				fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("  -%s %T", flagName(field.name), *f), "\n")
			case *string:
				fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("  -%s %T", flagName(field.name), *f), "\n")
			}
			fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("      %s", field.description), "\n")
		}
	}
}
