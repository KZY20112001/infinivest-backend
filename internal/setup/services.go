package setup

import (
	"github.com/KZY20112001/infinivest-backend/internal/cache"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
)

func InitUserService(db *gorm.DB) services.UserService {
	repo := repositories.NewPostgresUserRepo(db)
	return services.NewUserServiceImpl(repo)
}

func InitProfileService(db *gorm.DB, us services.UserService) services.ProfileService {
	repo := repositories.NewPostgresProfileRepo(db)
	return services.NewProfileServiceImpl(repo, us)
}

func InitS3Service(client *s3.PresignClient) services.S3Service {
	repo := repositories.NewS3RepositoryImpl(client)
	return services.NewS3ServiceImpl(repo)
}

func InitGenAIService() services.GenAIService {
	baseUrl := "http://localhost:5000"
	genAIRepo := repositories.NewFlaskMicroservice(baseUrl)
	return services.NewGenAIService(genAIRepo)
}

func InitRoboPortfolioService(pr repositories.PortfolioRepo, pc cache.PortfolioCache, ps services.ProfileService, gs services.GenAIService) services.RoboPortfolioService {
	return services.NewRoboPortfolioService(pr, pc, ps, gs)
}

func InitManualPortfolioService(pr repositories.PortfolioRepo, pc cache.PortfolioCache, ps services.ProfileService) services.ManualPortfolioService {
	return services.NewManualPortfolioService(pr, pc, ps)
}
