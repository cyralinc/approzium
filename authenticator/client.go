package main

import (
	"context"
	pb "dbauth/authenticator/messages"
	"fmt"

	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("enter client test\n")

	conn, err := grpc.Dial("authenticator:1234", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	fmt.Printf("dialed grpc server\n")

	client := pb.NewAuthenticatorClient(conn)
	ctx := context.Background()
	req := pb.AuthenticateRequest{
		Identity: "diotim",
		Salt:     []byte("swag"),
	}
	resp, err := client.Authenticate(ctx, &req)
	if err != nil {
		panic(fmt.Sprintf("failed to authenticate: %v", err))
	}
	fmt.Printf("made grpc request\n")

	creds := resp.GetCredentials()
	if creds == nil {
		panic("no credentials returned")
	}

	fmt.Printf("\nuser: %s\thashedPassword: %s\n", creds.User, creds.HashedPassword)
}
