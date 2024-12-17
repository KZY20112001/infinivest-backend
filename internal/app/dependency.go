package app

import (
	"github.com/KZY20112001/infinivest-backend/internal/cache"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

var (
	UserService    *services.UserService
	ProfileService *services.ProfileService
)

func InjectDependencies() {
	redisCache := cache.NewRedisCache(redisClient)
	postgresUserRepo := repositories.NewPostgresUserRepo(postgresDB)
	postgresProfileRepo := repositories.NewPostgresProfileRepo(postgresDB)

	UserService = services.NewUserService(postgresUserRepo, redisCache)
	ProfileService = services.NewProfileService(postgresProfileRepo, *UserService)
}
