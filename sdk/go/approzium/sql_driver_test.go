package approzium

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestAddPlaceholderPassword(t *testing.T) {
	if _, err := addPlaceholderPassword("postgres://pqgotest:somepass@localhost/pqgotest?sslmode=verify-full"); err == nil {
		t.Fatal("expected err because password should not be included in string")
	}

	if _, err := addPlaceholderPassword("user=postgres password=somepass dbname=postgres host=localhost port=5432 sslmode=disable"); err == nil {
		t.Fatal("expected err because password should not be included in string")
	}

	result, err := addPlaceholderPassword("postgres://pqgotest:@localhost/pqgotest?sslmode=verify-full")
	if err != nil {
		t.Fatal(err)
	}
	if result != "postgres://pqgotest:unknown@localhost/pqgotest?sslmode=verify-full" {
		t.Fatalf("expected postgres://pqgotest:unknown@localhost/pqgotest?sslmode=verify-full, but received %s", result)
	}

	result, err = addPlaceholderPassword("user=postgres dbname=postgres host=localhost port=5432 sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	if result != "user=postgres dbname=postgres host=localhost port=5432 sslmode=disable password=unknown" {
		t.Fatalf("expected user=postgres dbname=postgres host=localhost port=5432 sslmode=disable password=unknown but received %s", result)
	}
}

func TestParseDSN(t *testing.T) {
	logger := log.New()

	if _, _, err := parseDSN(logger, "postgres://pqgotest:@/pqgotest?sslmode=verify-full"); err == nil {
		t.Fatal("expected err because host must be included")
	}

	if _, _, err := parseDSN(logger, "user=postgres dbname=postgres port=5432 sslmode=disable"); err == nil {
		t.Fatal("expected err because host must be included")
	}

	dbHost, dbPort, err := parseDSN(logger, "postgres://pqgotest:@localhost/pqgotest?sslmode=verify-full")
	if err != nil {
		t.Fatal(err)
	}
	if dbHost != "localhost" {
		t.Fatalf("expected localhost, but received %s", dbHost)
	}
	if dbPort != defaultPostgresPort {
		t.Fatalf("expected %s, but received %s", defaultPostgresPort, dbPort)
	}

	dbHost, dbPort, err = parseDSN(logger, "user=postgres dbname=postgres host=localhost sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	if dbHost != "localhost" {
		t.Fatalf("expected localhost, but received %s", dbHost)
	}
	if dbPort != defaultPostgresPort {
		t.Fatalf("expected %s, but received %s", defaultPostgresPort, dbPort)
	}

	dbHost, dbPort, err = parseDSN(logger, "postgres://pqgotest:@localhost:1234/pqgotest?sslmode=verify-full")
	if err != nil {
		t.Fatal(err)
	}
	if dbHost != "localhost" {
		t.Fatalf("expected localhost, but received %s", dbHost)
	}
	if dbPort != "1234" {
		t.Fatalf("expected 1234, but received %s", dbPort)
	}

	dbHost, dbPort, err = parseDSN(logger, "user=postgres dbname=postgres host=localhost port=1234 sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	if dbHost != "localhost" {
		t.Fatalf("expected localhost, but received %s", dbHost)
	}
	if dbPort != "1234" {
		t.Fatalf("expected 1234, but received %s", dbPort)
	}
}
