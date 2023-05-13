package redis

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const defaultDB = 1

func OpenClient() *redis.Client {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		db = defaultDB
	}
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	_, err = client.Ping(context.Background()).Result()
	
	if err == nil {
		log.Printf("Redis connection to %s:%s is successfully connected!\n", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	} else {
		log.Fatalf("Redis connection to %s:%s is failed: %s",  os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), err.Error())
	}

	return client
}

