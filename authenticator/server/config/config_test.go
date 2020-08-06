package config

import (
	"flag"
	"os"
	"testing"
)

// This should only be set for local development. The tests behind this flag succeed when run
// one at a time, but fail when run all at once together because the flags and env variables
// they set interact with each other.
const envVarTestConfDir = "TEST_CONF_DIR"

func TestParseConfigFromEnv(t *testing.T) {
	os.Unsetenv("APPROZIUM_HOST")
	os.Unsetenv("APPROZIUM_HTTP_PORT")
	os.Unsetenv("APPROZIUM_LOG_LEVEL")
	os.Setenv("APPROZIUM_DISABLE_TLS", "true")
	config, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	if config.Host != "127.0.0.1" {
		t.Fatalf("expected %s, received %s", "127.0.0.1", config.Host)
	}
	if config.HTTPPort != 6000 {
		t.Fatalf("expected %d, received %d", 6000, config.HTTPPort)
	}
	if config.LogLevel != "info" {
		t.Fatalf("expected %s, received %s", "info", config.LogLevel)
	}

	os.Setenv("APPROZIUM_HOST", "0.0.0.0")
	os.Setenv("APPROZIUM_HTTP_PORT", "6001")
	os.Setenv("APPROZIUM_LOG_LEVEL", "debug")
	os.Setenv("APPROZIUM_DISABLE_TLS", "true")
	config, err = Parse()
	if err != nil {
		t.Fatal(err)
	}
	if config.Host != "0.0.0.0" {
		t.Fatalf("expected %s, received %s", "0.0.0.0", config.Host)
	}
	if config.HTTPPort != 6001 {
		t.Fatalf("expected %d, received %d", 6001, config.HTTPPort)
	}
	if config.LogLevel != "debug" {
		t.Fatalf("expected %s, received %s", "debug", config.LogLevel)
	}
	if !config.DisableTLS {
		t.Fatal("tls should be disabled")
	}
}

func TestParseConfigVersion(t *testing.T) {
	if err := flag.Set("disable-tls", "true"); err != nil {
		t.Fatal(err)
	}
	if err := flag.Set("version", "true"); err != nil {
		t.Fatal(err)
	}
	config, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	if !config.Version {
		t.Fatal("config version flag should be on")
	}
}

func TestParseConfigDevMode(t *testing.T) {
	if err := flag.Set("dev", "true"); err != nil {
		t.Fatal(err)
	}
	config, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	if !config.DevMode {
		t.Fatal("should be in dev mode")
	}
}

func TestPrecedenceFlagFirst(t *testing.T) {
	testConfDir := os.Getenv(envVarTestConfDir)
	if testConfDir == "" {
		t.Skip("skipping because TEST_CONF_DIR is unset")
	}

	if err := flag.Set("config", testConfDir+"/approzium.yaml"); err != nil {
		t.Fatal(err)
	}
	if err := flag.Set("http-port", "1"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("APPROZIUM_HTTP_PORT", "2"); err != nil {
		t.Fatal(err)
	}

	config, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	if config.HTTPPort != 1 {
		t.Fatalf("expected 1 but received %d, should be 1 from the flag because it should take precedence", config.HTTPPort)
	}
}

func TestPrecedenceEnvSecond(t *testing.T) {
	testConfDir := os.Getenv(envVarTestConfDir)
	if testConfDir == "" {
		t.Skip("skipping because TEST_CONF_DIR is unset")
	}

	if err := os.Setenv("APPROZIUM_HTTP_PORT", "2"); err != nil {
		t.Fatal(err)
	}

	if err := flag.Set("config", testConfDir+"/approzium.yaml"); err != nil {
		t.Fatal(err)
	}

	config, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	if config.HTTPPort != 2 {
		t.Fatalf("expected 1 but received %d, should be 1 from the env var because it should take precedence", config.HTTPPort)
	}
}

func TestPrecedenceFileThird(t *testing.T) {
	testConfDir := os.Getenv(envVarTestConfDir)
	if testConfDir == "" {
		t.Skip("skipping because TEST_CONF_DIR is unset")
	}

	if err := flag.Set("config", testConfDir+"/approzium.yaml"); err != nil {
		t.Fatal(err)
	}

	config, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	if config.HTTPPort != 3 {
		t.Fatalf("expected 3 but received %d, should be 3 from the config file because it should take precedence", config.HTTPPort)
	}
}

func TestPrecedenceDefaultLast(t *testing.T) {
	testConfDir := os.Getenv(envVarTestConfDir)
	if testConfDir == "" {
		t.Skip("skipping because TEST_CONF_DIR is unset")
	}

	if err := flag.Set("disable-tls", "true"); err != nil {
		t.Fatal(err)
	}
	config, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	if config.HTTPPort != 6000 {
		t.Fatalf("expected 6000 but received %d, should be 6000 from the default config because it should take precedence", config.HTTPPort)
	}
}
