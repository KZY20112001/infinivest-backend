package db

import (
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Redis *redis.Client

func ConnectToPostgres() {
	log.Println("Starting connection with Postgres Db")

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, db, port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	log.Println("Successfully connected to Postgres Db")

}

func ConnectToRedis() {
	log.Println("Starting connection with Redis Db")
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	address := fmt.Sprintf("%s:%s", host, port)
	Redis = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	log.Println("Successfully connected to Redis Db")
}
