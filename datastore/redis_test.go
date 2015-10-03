package datastore

import (
	"os"
	"testing"
)

func TestRedisAddress(t *testing.T) {
	oldVariable := os.Getenv(redisEnvironmentVariable)
	os.Setenv(redisEnvironmentVariable, "")
	if redisHost := redisAddress(); redisHost != redisDefaultHost {
		t.Error("Expected default host", redisDefaultHost, "but got", redisHost)
	}

	anotherValue := "something else"
	os.Setenv(redisEnvironmentVariable, anotherValue)
	if redisHost := redisAddress(); redisHost != anotherValue {
		t.Error("Expected env value", anotherValue, "but got", redisHost)
	}

	os.Setenv(redisEnvironmentVariable, oldVariable)
}

func TestConnect(t *testing.T) {
	if err := redisConnect(); err != nil {
		t.Error("Error connecting to redis:", err)
	}
}
