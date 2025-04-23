package services

import (
	"context"

	"github.com/KZY20112001/infinivest-backend/internal/redis"
)

type NotificationService interface {
	GetNotifications(ctx context.Context, userID, limit uint) ([]string, error)
}

type notificationServiceImpl struct {
	redis redis.NotificationRedis
}

func NewNotificationService(r redis.NotificationRedis) *notificationServiceImpl {
	return &notificationServiceImpl{redis: r}
}

func (ns *notificationServiceImpl) GetNotifications(ctx context.Context, userID, limit uint) ([]string, error) {
	return ns.redis.GetNotifications(ctx, userID, limit)
}
