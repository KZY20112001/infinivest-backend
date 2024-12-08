package db

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis() (*redis.Client, error) {
	log.Println("Connecting to Redis DB...")

	address := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	if client == nil {
		return nil, errors.New("failed to connect to Redis")
	}

	log.Println("Successfully connected to Redis DB")
	return client, nil
}
