package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/app"
	"github.com/KZY20112001/infinivest-backend/internal/db"
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/routes"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	postgresDB  *gorm.DB
	redisClient *redis.Client
)

func init() {
	app.LoadEnv()

	var err error
	postgresDB, err = db.ConnectToPostgres()
	if err != nil {
		log.Fatalf("error in connecting to database: %v", err.Error())
	}

	redisClient, err = db.ConnectToRedis()
	if err != nil {
		log.Fatalf("error in connecting to redis: %v", err.Error())
	}
}

func main() {
	r := gin.Default()

	redisCache := repositories.NewRedisCache(redisClient)
	postgresUserRepo := repositories.NewPostgresUserRepo(postgresDB)
	postgresProfileRepo := repositories.NewPostgresProfileRepo(postgresDB)

	userService := services.NewUserServiceImpl(postgresUserRepo, redisCache)
	profileService := services.NewProfileServiceImpl(postgresProfileRepo, userService)

	userHandler := handlers.NewUserHandlerImpl(userService)
	profileHandler := handlers.NewProfileHandlerImpl(profileService)

	routes.RegisterRoutes(r, userHandler, profileHandler)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
