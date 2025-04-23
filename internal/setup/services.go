package setup

import (
	"github.com/KZY20112001/infinivest-backend/internal/redis"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

func Services(
	portfolioRedis redis.RoboPortfolioRedis,
	notificationRedis redis.NotificationRedis,
	userRepo repositories.UserRepo,
	profileRepo repositories.ProfileRepo,
	roboPortfolioRepo repositories.RoboPortfolioRepo,
	manualPortfolioRepo repositories.ManualPortfolioRepo,
	s3Repo repositories.S3Repository,
	genAIRepo repositories.GenAIRepository,
) (
	services.UserService,
	services.ProfileService,
	services.RoboPortfolioService,
	services.ManualPortfolioService,
	services.NotificationService,
	services.S3Service,
	services.GenAIService,
) {
	userService := services.NewUserServiceImpl(userRepo)

	profileService := services.NewProfileServiceImpl(profileRepo, userService)

	s3Service := services.NewS3ServiceImpl(s3Repo)

	genAIService := services.NewGenAIService(genAIRepo)

	notificationService := services.NewNotificationService(notificationRedis)
	roboPortfolioService := services.NewRoboPortfolioService(
		roboPortfolioRepo, portfolioRedis, genAIService, notificationService, userService,
	)

	manualPortfolioService := services.NewManualPortfolioService(
		manualPortfolioRepo, genAIService,
	)

	return userService, profileService, roboPortfolioService, manualPortfolioService, notificationService, s3Service, genAIService
}
