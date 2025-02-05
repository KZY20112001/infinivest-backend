package main

import (
	"context"
	"fmt"
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
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	// 	log.Fatalf("error in connecting to redis: %w", err.Error())
	// }
}

func initUserService(db *gorm.DB) services.UserService {
	repo := repositories.NewPostgresUserRepo(db)
	return services.NewUserServiceImpl(repo)
}

func initProfileService(db *gorm.DB, userService services.UserService) services.ProfileService {
	repo := repositories.NewPostgresProfileRepo(db)
	return services.NewProfileServiceImpl(repo, userService)
}

func initS3Service(client *s3.PresignClient) services.S3Service {
	repo := repositories.NewS3RepositoryImpl(client)
	return services.NewS3ServiceImpl(repo)
}

func initPortfolioService() services.PortfolioService {
	baseUrl := "http://localhost:5000"
	repo := repositories.NewFlaskMicroservice(baseUrl)
	return services.NewPortfolioServiceImpl(repo)
}

func initHandlers(db *gorm.DB, s3Client *s3.PresignClient) (*handlers.UserHandler, *handlers.ProfileHandler, *handlers.PortfolioHandler, *handlers.S3Handler) {
	s3Service := initS3Service(s3Client)
	userService := initUserService(db)
	profileService := initProfileService(db, userService)
	portfolioService := initPortfolioService()
	return handlers.NewUserHandler(userService),
		handlers.NewProfileHandler(profileService),
		handlers.NewPortfolioHandler(portfolioService),
		handlers.NewS3Handler(s3Service)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)
	userHandler, profileHandler, portfolioHandler, s3Handler := initHandlers(postgresDB, presignClient)

	r := routes.RegisterRoutes(userHandler, profileHandler, portfolioHandler, s3Handler)
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
		fmt.Println("Server started")
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
