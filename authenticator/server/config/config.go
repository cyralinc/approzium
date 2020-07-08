package config

import (
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

	LogLevel       string
	LogFormat      string
	LogRaw         bool
	VaultTokenPath string
	ConfigFilePath string
}

// ParseConfig returns the parsed config. A pointer is not returned
// because after first parse, the config is immutable.
func ParseConfig() (Config, error) {
	var config Config
	setConfigDefaults()
	setConfigEnvVars()
	err := viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func setConfigEnvVars() {
	viper.SetEnvPrefix("approzium")
	viper.BindEnv("Host", "APPROZIUM_HOST")
	viper.BindEnv("HTTPPort", "APPROZIUM_HTTP_PORT")
	viper.BindEnv("GRPCPort", "APPROZIUM_GRPC_PORT")
	viper.BindEnv("LogLevel", "APPROZIUM_LOG_LEVEL")
	viper.BindEnv("LogFormat", "APPROZIUM_LOG_FORMAT")
	viper.BindEnv("LogRaw", "APPROZIUM_LOG_RAW")
	viper.BindEnv("VaultTokenPath", "APPROZIUM_VAULT_TOKEN_PATH")
	viper.BindEnv("ConfigFilePath", "APPROZIUM_CONFIG_FILE_PATH")
}

func setConfigDefaults() {
	viper.SetDefault("Host", "127.0.0.1")
	viper.SetDefault("HTTPPort", "6000")
	viper.SetDefault("GRPCPort", "6001")
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("LogFormat", "text")
	viper.SetDefault("LogRaw", false)
}
