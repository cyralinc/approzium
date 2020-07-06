package config

import "github.com/kelseyhightower/envconfig"

// Config is an object for storing configuration variables set through
// Approzium's environment.
//
// Please see https://approzium.org/configuration for elaboration upon
// each parameter.
type Config struct {
	Host     string `default:"127.0.0.1"`
	HTTPPort int    `envconfig:"http_port" default:"6000"`
	GRPCPort int    `envconfig:"grpc_port" default:"6001"`

	LogLevel  string `envconfig:"log_level" default:"info"`
	LogFormat string `envconfig:"log_format" default:"text"` // Also supports "json".
	LogRaw    bool   `envconfig:"log_raw" default:"false"`

	VaultTokenPath string `envconfig:"vault_token_path"`
}

// ParseConfig returns the parsed config. A pointer is not returned
// because after first parse, the config is immutable.
func ParseConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("approzium", &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
