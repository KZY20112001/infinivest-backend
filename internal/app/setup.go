package app

import (
	"errors"
	"log"
	"sync"

	"github.com/KZY20112001/infinivest-backend/internal/db"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func LoadEnv() {
	err := godotenv.Load(".env.local")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func SetupDB() {
	log.Println("Connecting to Postgres and Redis...")

	var wg sync.WaitGroup
	wg.Add(2)

	errCh := make(chan error, 2)

	// Connect to Postgres
	go func() {
		defer wg.Done()
		log.Println("Connecting to Postgres DB...")
		var err error
		DB, err = db.ConnectToPostgres()
		if err != nil {
			errCh <- err
			return
		}
		log.Println("Successfully connected to Postgres DB")
	}()

	// Connect to Redis
	go func() {
		defer wg.Done()
		log.Println("Connecting to Redis...")
		Redis = db.ConnectToRedis() // Assuming this function already handles errors internally
		if Redis == nil {
			errCh <- errors.New("failed to connect to Redis")
			return
		}
		log.Println("Successfully connected to Redis")
	}()

	// Wait for both tasks to complete
	wg.Wait()
	close(errCh)

	// Check for any errors
	for err := range errCh {
		log.Fatal(err.Error()) // Log and terminate if any error occurs
	}

	log.Println("Connected to Postgres and Redis successfully")

	DB.AutoMigrate(&models.User{}, &models.Profile{})
}
