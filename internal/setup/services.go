package setup

import (
	"github.com/KZY20112001/infinivest-backend/internal/caches"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

func Services(
	portfolioCache caches.RoboPortfolioCache,
	userRepo repositories.UserRepo,
	profileRepo repositories.ProfileRepo,
	portfolioRepo repositories.RoboPortfolioRepo,
	s3Repo repositories.S3Repository,
	genAIRepo repositories.GenAIRepository,
) (
	services.UserService,
	services.ProfileService,
	services.RoboPortfolioService,
	services.ManualPortfolioService,
	services.S3Service,
	services.GenAIService,
) {
	userService := services.NewUserServiceImpl(userRepo)

	profileService := services.NewProfileServiceImpl(profileRepo, userService)

	s3Service := services.NewS3ServiceImpl(s3Repo)

	genAIService := services.NewGenAIService(genAIRepo)

	roboPortfolioService := services.NewRoboPortfolioService(
		portfolioRepo, portfolioCache, profileService, genAIService,
	)

	manualPortfolioService := services.NewManualPortfolioService(
		portfolioRepo, portfolioCache, profileService,
	)

	return userService, profileService, roboPortfolioService, manualPortfolioService, s3Service, genAIService
}
