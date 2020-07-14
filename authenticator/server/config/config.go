package config

import (
	"errors"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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

	VaultTokenPath string
	ConfigFilePath string

	// Special flags
	devMode bool
	Version bool
}

// Parse returns the parsed config. A pointer is not returned
// because after first parse, the config is immutable.
func Parse() (Config, error) {
	config, err := parse()
	if err != nil {
		return Config{}, err
	}
	if err := verify(config); err != nil {
		return Config{}, err
	}
	if config.devMode {
		// dev mode uses the file back-end
		os.Unsetenv("VAULT_ADDR")
		config = Config{
			Host:       "127.0.0.1",
			HTTPPort:   6000,
			GRPCPort:   6001,
			DisableTLS: true,
			LogLevel:   "debug",
			LogFormat:  "text",
			LogRaw:     false,
		}
	}
	return config, nil
}

func parse() (Config, error) {
	var config Config
	setConfigDefaults()
	setConfigFlags()
	if err := setConfigEnvVars(); err != nil {
		return Config{}, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	if config.ConfigFilePath == "" {
		return config, nil
	}
	viper.SetConfigName("approzium.config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(config.ConfigFilePath)

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func verify(config Config) error {
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

func setConfigEnvVars() error {
	viper.SetEnvPrefix("approzium")
	if err := viper.BindEnv("Host", "APPROZIUM_HOST"); err != nil {
		return err
	}
	if err := viper.BindEnv("HTTPPort", "APPROZIUM_HTTP_PORT"); err != nil {
		return err
	}
	if err := viper.BindEnv("GRPCPort", "APPROZIUM_GRPC_PORT"); err != nil {
		return err
	}

	if err := viper.BindEnv("DisableTLS", "APPROZIUM_DISABLE_TLS"); err != nil {
		return err
	}
	if err := viper.BindEnv("PathToTLSCert", "APPROZIUM_PATH_TO_TLS_CERT"); err != nil {
		return err
	}
	if err := viper.BindEnv("PathToTLSKey", "APPROZIUM_PATH_TO_TLS_KEY"); err != nil {
		return err
	}

	if err := viper.BindEnv("LogLevel", "APPROZIUM_LOG_LEVEL"); err != nil {
		return err
	}
	if err := viper.BindEnv("LogFormat", "APPROZIUM_LOG_FORMAT"); err != nil {
		return err
	}
	if err := viper.BindEnv("LogRaw", "APPROZIUM_LOG_RAW"); err != nil {
		return err
	}

	if err := viper.BindEnv("VaultTokenPath", "APPROZIUM_VAULT_TOKEN_PATH"); err != nil {
		return err
	}
	if err := viper.BindEnv("ConfigFilePath", "APPROZIUM_CONFIG_FILE_PATH"); err != nil {
		return err
	}
	return nil
}

func setConfigDefaults() {
	viper.SetDefault("Host", "127.0.0.1")
	viper.SetDefault("HTTPPort", "6000")
	viper.SetDefault("GRPCPort", "6001")

	viper.SetDefault("DisableTLS", false)

	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("LogFormat", "text")
	viper.SetDefault("LogRaw", false)
}

func setConfigFlags() {
	if pflag.Lookup("host") == nil {
		// avoid redefining flags because it leads to panic
		pflag.String("host", "", "")
		pflag.String("httpport", "", "")
		pflag.String("grpcport", "", "")

		pflag.Bool("disabletls", false, "")
		pflag.String("tlscertpath", "", "")
		pflag.String("tlskeypath", "", "")

		pflag.String("loglevel", "", "")
		pflag.String("logformat", "", "")
		pflag.Bool("lograw", false, "")

		pflag.String("vaulttokenpath", "", "")
		pflag.String("config", "", "")

		pflag.Bool("dev", false, "whether to run the authenticator in dev mode")
		pflag.Bool("version", false, "output the current authenticator version")
	}

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.BindPFlag("configfilepath", pflag.Lookup("config"))
}
