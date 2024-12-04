package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Postgres driver
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ConnectToPostgres connects to the PostgreSQL database
func ConnectToPostgres() (*sql.DB, error) {
	fmt.Println("Starting connection with Postgres Db")
	db, err := sql.Open("postgres", "user=admin password=admin dbname=postgres_db sslmode=disable port=5432")
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err = db.Ping(); err != nil {
		log.Println("DB Ping Failed")
		return nil, err
	}

	log.Println("DB Connection started successfully")
	return db, nil
}

// ConnectToRedis connects to the Redis server
func ConnectToRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	// Check connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Println("Redis Ping Failed")
		return nil, err
	}

	log.Println("Redis Connection started successfully")
	return client, nil
}
