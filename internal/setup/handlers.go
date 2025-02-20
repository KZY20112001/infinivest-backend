package setup

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

func Handlers(
	userService services.UserService,
	profileService services.ProfileService,
	roboPortfolioService services.RoboPortfolioService,
	manualPortfolioService services.ManualPortfolioService,
	s3Service services.S3Service,
	genAIService services.GenAIService,
) (
	*handlers.UserHandler,
	*handlers.ProfileHandler,
	*handlers.RoboPortfolioHandler,
	*handlers.ManualPortfolioHandler,
	*handlers.S3Handler,
) {
	return handlers.NewUserHandler(userService),
		handlers.NewProfileHandler(profileService),
		handlers.NewRoboPortfolioHandler(roboPortfolioService, genAIService),
		handlers.NewManualPortfolioHandler(manualPortfolioService),
		handlers.NewS3Handler(s3Service)
}
