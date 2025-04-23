package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type NotificationRedis interface {
	AddNotification(ctx context.Context, userID uint, notification string) error
	GetNotifications(ctx context.Context, userID, limit uint) ([]string, error)
}

type notificationRedis struct {
	client *redis.Client
}

func NewNotificationRedis(client *redis.Client) *notificationRedis {
	return &notificationRedis{client: client}
}

func (r *notificationRedis) AddNotification(ctx context.Context, userID uint, notification string) error {
	key := fmt.Sprintf("notification_queue:%d", userID)
	return r.client.LPush(ctx, key, notification).Err()
}

func (r *notificationRedis) GetNotifications(ctx context.Context, userID, limit uint) ([]string, error) {
	key := fmt.Sprintf("notification_queue:%d", userID)

	var end int64 = 0
	if limit > 0 {
		end = int64(limit) - 1
	}

	return r.client.LRange(ctx, key, 0, end).Result()
}
