package setup

import (
	customRedis "github.com/KZY20112001/infinivest-backend/internal/redis"
	"github.com/redis/go-redis/v9"
)

func Redis(client *redis.Client) (customRedis.RoboPortfolioRedis, customRedis.NotificationRedis) {
	return customRedis.NewRoboPortfolioRedis(client), customRedis.NewNotificationRedis(client)
}
