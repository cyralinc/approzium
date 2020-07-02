package server

import "github.com/kelseyhightower/envconfig"

// Config is an object for storing configuration variables set through
// Approzium's environment. Supports:
// 		- APPROZIUM_HOST: Defaults to 127.0.0.1.
//		- APPROZIUM_PORT: Defaults to 6000.
//		- APPROZIUM_LOG_LEVEL: Defaults to info. Valid values are:
//			- trace
//			- debug
//			- info
//			- warn
//			- error
//			- fatal
//			- panic
//
// For those using Vault for storage, Approzium will read Vault's address
// and the Vault token it should use through Vault's normal environment
// variables described here:
// https://www.vaultproject.io/docs/commands#environment-variables.
// At a minimum, VAULT_ADDR and VAULT_TOKEN must be set.
type Config struct {
	Host           string `default:"127.0.0.1"`
	Port           int    `default:"6000"`
	LogLevel       string `envconfig:"log_level" default:"info"`
	LogFormat      string `envconfig:"log_format" default:"text"` // Also supports "json".
	LogRaw         bool   `envconfig:"log_raw" default:"false"`
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
