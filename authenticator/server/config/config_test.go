package config

import (
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	os.Unsetenv("APPROZIUM_HOST")
	os.Unsetenv("APPROZIUM_PORT")
	os.Unsetenv("APPROZIUM_LOG_LEVEL")
	os.Setenv("APPROZIUM_DISABLE_TLS", "true")
	config, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	if config.Host != "127.0.0.1" {
		t.Fatalf("expected %s, received %s", "127.0.0.1", config.Host)
	}
	if config.Port != 6000 {
		t.Fatalf("expected %d, received %d", 6000, config.Port)
	}
	if config.LogLevel != "info" {
		t.Fatalf("expected %s, received %s", "info", config.LogLevel)
	}

	os.Setenv("APPROZIUM_HOST", "0.0.0.0")
	os.Setenv("APPROZIUM_PORT", "6001")
	os.Setenv("APPROZIUM_LOG_LEVEL", "debug")
	os.Setenv("APPROZIUM_DISABLE_TLS", "true")
	config, err = Parse()
	if err != nil {
		t.Fatal(err)
	}
	if config.Host != "0.0.0.0" {
		t.Fatalf("expected %s, received %s", "0.0.0.0", config.Host)
	}
	if config.Port != 6001 {
		t.Fatalf("expected %d, received %d", 6001, config.Port)
	}
	if config.LogLevel != "debug" {
		t.Fatalf("expected %s, received %s", "debug", config.LogLevel)
	}
	if !config.DisableTLS {
		t.Fatal("tls should be disabled")
	}
}
