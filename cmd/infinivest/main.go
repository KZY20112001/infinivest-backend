package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/conf"
	"github.com/KZY20112001/infinivest-backend/internal/db"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/routes"
	"github.com/KZY20112001/infinivest-backend/internal/setup"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	postgresDB  *gorm.DB
	redisClient *redis.Client
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
	postgresDB.AutoMigrate(&models.User{}, &models.Profile{}, &models.Portfolio{}, &models.PortfolioCategory{}, &models.PortfolioAsset{})

	redisClient, err = db.ConnectToRedis()
	if err != nil {
		log.Fatalf("error in connecting to redis: %v", err.Error())
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	appConf := conf.LoadConfig()

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	// init repositories
	userRepo, profileRepo, portfolioRepo, s3Repo, genAIRepo := setup.Repositories(
		postgresDB, presignClient, appConf.FlaskMicroserviceURL,
	)

	// init caches
	portfolioCache := setup.Caches(redisClient)

	// init services
	userService, profileService, roboPortfolioService, manualPortfolioService, s3Service, genAIService := setup.Services(
		portfolioCache, userRepo, profileRepo, portfolioRepo, s3Repo, genAIRepo,
	)

	// init handlers
	userHandler, profileHandler, roboPortfolioHandler, manualPortfolioHandler, s3Handler := setup.Handlers(
		userService, profileService, roboPortfolioService, manualPortfolioService, s3Service, genAIService,
	)

	// init schedulers
	portfolioScheduler := setup.PortfolioScheduler(
		roboPortfolioService, portfolioRepo, portfolioCache,
	)

	portfolioScheduler.Start(ctx)
	r := routes.RegisterRoutes(userHandler, profileHandler, roboPortfolioHandler, manualPortfolioHandler, s3Handler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
		log.Println("Server started")
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	stop()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
