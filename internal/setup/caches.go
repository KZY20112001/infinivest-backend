package setup

import (
	"github.com/KZY20112001/infinivest-backend/internal/caches"
	"github.com/redis/go-redis/v9"
)

func Caches(client *redis.Client) caches.PortfolioCache {
	return caches.NewPortfolioRedis(client)
}
