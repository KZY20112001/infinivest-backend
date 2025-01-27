package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/db"
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/routes"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	postgresDB *gorm.DB
	// redisClient *redis.Client
)

func init() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgresDB, err = db.ConnectToPostgres()
	if err != nil {
		log.Fatalf("error in connecting to database: %v", err.Error())
	}
	postgresDB.AutoMigrate(&models.User{}, &models.Profile{})

	// redisClient, err = db.ConnectToRedis()
	// if err != nil {
	// 	log.Fatalf("error in connecting to redis: %v", err.Error())
	// }
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// redisCache := repositories.NewRedisCache(redisClient)
	postgresUserRepo := repositories.NewPostgresUserRepo(postgresDB)
	postgresProfileRepo := repositories.NewPostgresProfileRepo(postgresDB)

	userService := services.NewUserServiceImpl(postgresUserRepo)
	profileService := services.NewProfileServiceImpl(postgresProfileRepo, userService)

	userHandler := handlers.NewUserHandler(userService)
	profileHandler := handlers.NewProfileHandler(profileService)

	r := routes.RegisterRoutes(userHandler, profileHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
