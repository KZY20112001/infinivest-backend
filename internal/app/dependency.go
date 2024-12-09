package app

import (
	"github.com/KZY20112001/infinivest-backend/internal/cache"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

var UserService *services.UserService

func InjectDependencies() {
	postgresUserRepo := repositories.NewPostgresUserRepo(postgresDB)
	redisCache := cache.NewRedisCache(redisClient)
	UserService = services.NewUserService(postgresUserRepo, redisCache)
}
