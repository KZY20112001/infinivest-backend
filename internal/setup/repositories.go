package setup

import (
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
)

func Repositories(
	db *gorm.DB,
	s3Client *s3.PresignClient,
	genAIUrl string,
) (
	repositories.UserRepo,
	repositories.ProfileRepo,
	repositories.RoboPortfolioRepo,
	repositories.ManualPortfolioRepo,
	repositories.S3Repository,
	repositories.GenAIRepository,
) {
	return repositories.NewPostgresUserRepo(db),
		repositories.NewPostgresProfileRepo(db),
		repositories.NewPostgresRoboPortfolioRepo(db),
		repositories.NewPostgresManualPortfolioRepo(db),
		repositories.NewS3RepositoryImpl(s3Client),
		repositories.NewFlaskMicroservice(genAIUrl)
}
