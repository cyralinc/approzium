package main

import (
	pb "github.com/approzium/approzium/authenticator/protos"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

const (
	serviceAddress = ":1234"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	authenticator, err := NewAuthenticator()
	if err != nil {
		log.Panicf("failed to create authenticator: %s", err)
	}

	lis, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		log.Panicf("failed to listen on %s", serviceAddress)
	}
	log.Infof("authenticator listening for requests on %s\n", serviceAddress)

	go authenticator.run()

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticatorServer(grpcServer, authenticator)
	if err := grpcServer.Serve(lis); err != nil {
		log.Panicf("failed to listen on %s", serviceAddress)
	}
}
