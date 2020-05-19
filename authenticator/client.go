package main

import (
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("authenticator:1234", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	fmt.Printf("connected\n")
	defer conn.Close()
}
