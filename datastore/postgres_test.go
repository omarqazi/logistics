package datastore

import (
	"os"
	"testing"
)

func TestPostgresAddress(t *testing.T) {
	oldVariable := os.Getenv(postgresEnvironmentVariable)
	os.Setenv(postgresEnvironmentVariable, "")
	if postgresHost := postgresAddress(); postgresHost != postgresDefaultAddress {
		t.Error("Expected default host", postgresDefaultAddress, "but got", postgresHost)
	}

	anotherValue := "something else"
	os.Setenv(postgresEnvironmentVariable, anotherValue)
	if postgresHost := postgresAddress(); postgresHost != anotherValue {
		t.Error("Expected env value", anotherValue, "but got", postgresHost)
	}

	os.Setenv(postgresEnvironmentVariable, oldVariable)
}

func TestPostgresConnect(t *testing.T) {
	if err := postgresConnect(); err != nil {
		t.Error("Error connecting to redis:", err)
	}
}
