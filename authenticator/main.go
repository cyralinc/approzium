package main

import (
	"fmt"
	"net"
	"strings"

	pb "github.com/approzium/approzium/authenticator/protos"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

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
	Host     string `default:"127.0.0.1"`
	Port     int    `default:"6000"`
	LogLevel string `envconfig:"log_level" default:"info"`
}

func main() {
	config, err := parseConfig()
	if err != nil {
		log.Panicf("couldn't parse config: %s", err)
	}

	logLevel, err := log.ParseLevel(strings.ToLower(config.LogLevel))
	if err != nil {
		log.Panicf("couldn't parse log level: %s", err)
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})

	authenticator, err := NewAuthenticator()
	if err != nil {
		log.Panicf("failed to create authenticator: %s", err)
	}

	serviceAddress := fmt.Sprintf("%s:%d", config.Host, config.Port)
	lis, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		log.Panicf("failed to listen on %s: %s", serviceAddress, err)
	}
	log.Infof("authenticator listening for requests on %s\n", serviceAddress)

	go authenticator.run()

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticatorServer(grpcServer, authenticator)
	if err := grpcServer.Serve(lis); err != nil {
		log.Panicf("failed to listen on %s", serviceAddress)
	}
}

// parseConfig returns the parsed config. A pointer is not returned
// because after first parse, the config is immutable.
func parseConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("approzium", &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
