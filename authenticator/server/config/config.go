package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

// Config is an object for storing configuration variables set through
// Approzium's environment.
//
// Please see https://approzium.org/configuration for elaboration upon
// each parameter.
type Config struct {
	Host     string `yaml:"host" env:"APPROZIUM_HOST" env-default:"127.0.0.1"`
	HTTPPort int    `yaml:"http_port" env:"APPROZIUM_HTTP_PORT" env-default:"6000"`
	GRPCPort int    `yaml:"grpc_port" env:"APPROZIUM_GRPC_PORT" env-default:"6001"`

	LogLevel  string `yaml:"log_level" env:"APPROZIUM_LOG_LEVEL" env-default:"info"`
	LogFormat string `yaml:"log_format" env:"APPROZIUM_LOG_FORMAT" env-default:"text"` // Also supports "json".
	LogRaw    bool   `yaml:"log_raw" env:"APPROZIUM_LOG_RAW" env-default:"false"`

	VaultTokenPath string `yaml:"vault_token_path" env:"APPROZIUM_VAULT_TOKEN_PATH"`
	ConfigFilePath string `env:"APPROZIUM_CONFIG_FILE_PATH"`
}

// ParseConfig returns the parsed config. A pointer is not returned
// because after first parse, the config is immutable.
func ParseConfig() (Config, error) {
	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		return Config{}, err
	}
	if config.ConfigFilePath == "" {
		// if not config file path is provided and confif.yml exists in the current
		// directory, then use it
		if _, err := os.Stat("config.yml"); err == nil {
			config.ConfigFilePath = "config.yml"
		}
	}

	if config.ConfigFilePath != "" {
		err := cleanenv.ReadConfig(config.ConfigFilePath, &config)
		if err != nil {
			return Config{}, err
		}
	}
	return config, nil
}
