package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cyralinc/approzium/sdk/go/approzium"
)

func main() {
	// Create a connection to the Approzium authenticator,
	// because only the authenticator knows the password.
	authClient, err := approzium.NewAuthClient("localhost:6001", &approzium.Config{
		DisableTLS:      true,
		RoleArnToAssume: os.Getenv("TEST_ASSUMABLE_ARN"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Now create a Postgres connection without a password.
	// We also support strings like:
	// "postgres://pqgotest:@localhost/pqgotest?sslmode=verify-full"
	dataSourceName := "user=postgres dbname=postgres host=localhost port=5432 sslmode=disable"
	db, err := authClient.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		fmt.Println("successfully got a result")
	} else {
		fmt.Println("received nothing")
	}
}
