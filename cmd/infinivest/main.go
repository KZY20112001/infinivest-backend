package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/db"
	"github.com/KZY20112001/infinivest-backend/internal/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	var dbPostgres *sql.DB
	var dbRedis *redis.Client
	var err error

	go func() {
		// Initialize PostgreSQL connection
		dbPostgres, err = db.ConnectToPostgres()
		if err != nil {
			log.Fatalf("Error connecting to Postgres: %v", err)
		}
	}()

	go func() {
		// Initialize Redis connection
		dbRedis, err = db.ConnectToRedis()
		if err != nil {
			log.Fatalf("Error connecting to Redis: %v", err)
		}
	}()

	// Example: Wait for connections to be established before starting the server
	for err != nil {
		log.Println("Waiting for DB and Redis connections to be ready...")
		// Recheck connections
		if dbPostgres == nil {
			dbPostgres, err = db.ConnectToPostgres()
		}
		if dbRedis == nil {
			dbRedis, err = db.ConnectToRedis()
		}
		time.Sleep(2 * time.Second) // Retry interval
	}

	// Continue with the Gin server setup
	fmt.Println("Starting Gin server...")
	log.Println("DB Connection started successfully")
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run()
}
