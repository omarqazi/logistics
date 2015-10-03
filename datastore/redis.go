package datastore

import (
	"gopkg.in/redis.v3"
	"log"
	"os"
)

var Redis *redis.Client

const redisEnvironmentVariable = "REDIS_ADDRESS"
const redisDefaultHost = "localhost:6379"

func init() {
	if err := redisConnect(); err != nil {
		log.Fatalln(err)
	}
}

func redisConnect() error {
	Redis = redis.NewClient(&redis.Options{
		Addr: redisAddress(),
		DB:   0,
	})

	return Redis.Ping().Err()
}

func redisAddress() string {
	if raddr := os.Getenv(redisEnvironmentVariable); raddr != "" {
		return raddr
	}
	return redisDefaultHost
}
