package services

import (
	"context"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/redis"
)

type NotificationService interface {
	GetNotifications(ctx context.Context, userID, limit uint) ([]string, error)
	AddNotification(ctx context.Context, userID uint, notiType, message string) error
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

func (ns *notificationServiceImpl) AddNotification(ctx context.Context, userID uint, notiType, message string) error {
	time := time.Now().Format("2006-01-02 15:04:05")
	notification := time + "; " + notiType + "; " + message
	return ns.redis.AddNotification(ctx, userID, notification)
}
