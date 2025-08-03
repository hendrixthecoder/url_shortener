package main

import (
	"log"
	"os"

	"github.com/boj/redistore"
)

var store *redistore.RediStore

func InitRedisStore() {
	secret := os.Getenv("REDIS_SECRET")
	if secret == "" {
		log.Fatal("REDIS SECRET not provided")
	}

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Fatal("REDIS URL not provided")
	}

	var err error
	store, err = redistore.NewRediStore(10, "tcp", redisUrl, "", "", []byte(secret))
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
