package setup

import (
	"github.com/KZY20112001/infinivest-backend/internal/cache"
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitHandlers(db *gorm.DB, redisClient *redis.Client, s3Client *s3.PresignClient) (*handlers.UserHandler, *handlers.ProfileHandler, *handlers.RoboPortfolioHandler, *handlers.ManualPortfolioHandler, *handlers.S3Handler) {
	genAIService := InitGenAIService()
	s3Service := InitS3Service(s3Client)
	userService := InitUserService(db)
	profileService := InitProfileService(db, userService)

	portfolioRepo := repositories.NewPostgresPortfolioRepo(db)
	portfolioCache := cache.NewPortfolioRedis(redisClient)

	roboPortfolioService := InitRoboPortfolioService(portfolioRepo, portfolioCache, profileService, genAIService)
	manualPortfolioService := InitManualPortfolioService(portfolioRepo, portfolioCache, profileService)

	return handlers.NewUserHandler(userService),
		handlers.NewProfileHandler(profileService),
		handlers.NewRoboPortfolioHandler(roboPortfolioService, genAIService),
		handlers.NewManualPortfolioHandler(manualPortfolioService),
		handlers.NewS3Handler(s3Service)
}
