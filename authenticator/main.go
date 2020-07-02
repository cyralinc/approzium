package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/approzium/approzium/authenticator/server"
	pb "github.com/approzium/approzium/authenticator/server/protos"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	config, err := server.ParseConfig()
	if err != nil {
		log.Panicf("couldn't parse config: %s", err)
	}

	logLevel, err := log.ParseLevel(strings.ToLower(config.LogLevel))
	if err != nil {
		log.Panicf("couldn't parse log level: %s", err)
	}
	logger := log.New()
	logger.Level = logLevel

	switch strings.ToLower(config.LogFormat) {
	case "text":
		logger.SetFormatter(&log.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		})
	case "json":
		logger.SetFormatter(&log.JSONFormatter{})
	default:
		log.Panicf("unsupported log format: %s", config.LogFormat)
	}

	svr, err := server.New(logger, config)
	if err != nil {
		log.Panicf("failed to create authenticator: %s", err)
	}

	serviceAddress := fmt.Sprintf("%s:%d", config.Host, config.Port)
	lis, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		log.Panicf("failed to listen on %s: %s", serviceAddress, err)
	}
	log.Infof("authenticator listening for requests on %s\n", serviceAddress)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticatorServer(grpcServer, svr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Panicf("failed to listen on %s", serviceAddress)
	}
}
