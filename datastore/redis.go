package datastore

import (
	"gopkg.in/redis.v3"
	"log"
	"os"
)

var Redis *redis.Client

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr: redisAddress(),
		DB:   0,
	})

	err := Redis.Ping().Err()
	if err != nil {
		log.Fatalln(err)
	}
}

func redisAddress() string {
	if raddr := os.Getenv("REDIS_ADDRESS"); raddr != "" {
		return raddr
	}
	return "localhost:6379"
}
