package main

import (
	"fmt"
	"net"

	pb "dbauth/authenticator/messages"

	"google.golang.org/grpc"
)

const (
	serviceAddress = "localhost:1234"
)

func main() {
	authenticator := NewAuthenticator()

	lis, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		panic(fmt.Errorf("failed to listen on %s", serviceAddress))
	}
	fmt.Printf("authenticator listening for requests on %s", serviceAddress)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticatorServer(grpcServer, authenticator)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
