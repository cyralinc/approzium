package main

import (
	"context"
	pb "dbauth/authenticator/protos"
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
	req := pb.DBUserRequest{
		Identity: "diotim",
	}
	resp, err := client.GetDBUser(ctx, &req)
	if err != nil {
		panic(fmt.Sprintf("failed to execute GetDBUser: %v", err))
	}
	fmt.Printf("made GetDBUser request\n")

	user := resp.GetDbuser()
	if user == "" {
		panic("no dbuser received")
	}

	fmt.Printf("\nuser: %s\n", resp.Dbuser)

    req1 := pb.DBHashRequest{
		Identity: "diotim",
        Salt: []byte{0, 1, 0, 1},
	}
    resp1, err := client.GetDBHash(ctx, &req1)
	if err != nil {
		panic(fmt.Sprintf("failed to execute GetHashUser: %v", err))
	}
	fmt.Printf("made GetDBHash request\n")

	hash := resp1.GetHash()
	if hash == "" {
		panic("no hash received")
	}

	fmt.Printf("hash: %s\n", hash)
}
