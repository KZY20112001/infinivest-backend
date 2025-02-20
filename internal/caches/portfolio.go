package caches

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type PortfolioCache interface {
	AddPortfolioToRebalancingQueue(ctx context.Context, portfolioID uint, userID uint, nextRebalance time.Time) error
	GetNextRebalanceTime(ctx context.Context, portfolioID uint, userID uint) (time.Time, error)
	GetEarliestPortfolio(ctx context.Context) (uint, uint, time.Time, error)
}

type portfolioRedis struct {
	client *redis.Client
}

func NewPortfolioRedis(client *redis.Client) *portfolioRedis {
	return &portfolioRedis{client: client}
}

func (r *portfolioRedis) AddPortfolioToRebalancingQueue(ctx context.Context, portfolioID uint, userID uint, nextRebalance time.Time) error {
	key := "rebalancing_queue"
	score := float64(nextRebalance.Unix())
	member := fmt.Sprintf("%d:%d", userID, portfolioID)

	value := redis.Z{
		Score:  score,
		Member: member,
	}

	_, err := r.client.ZAdd(ctx, key, value).Result()
	return err
}

func (r *portfolioRedis) GetNextRebalanceTime(ctx context.Context, portfolioID uint, userID uint) (time.Time, error) {
	key := "rebalancing_queue"
	member := fmt.Sprintf("%d:%d", userID, portfolioID)

	nextRebalanceTime, err := r.client.ZScore(ctx, key, member).Result()

	if err != nil {
		if err == redis.Nil {
			return time.Time{}, fmt.Errorf("portfolio %d for user %d not found in rebalancing queue", portfolioID, userID)
		}
		return time.Time{}, err
	}
	return time.Unix(int64(nextRebalanceTime), 0), nil
}

func (r *portfolioRedis) GetEarliestPortfolio(ctx context.Context) (uint, uint, time.Time, error) {
	key := "rebalancing_queue"
	result, err := r.client.ZRangeWithScores(ctx, key, 0, 0).Result()
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	if len(result) == 0 {
		return 0, 0, time.Time{}, fmt.Errorf("no portfolios in the queue")
	}

	member := result[0].Member.(string) // "userID:portfolioID" string
	score := result[0].Score            // nextRebalance time

	var userID, portfolioID uint
	_, err = fmt.Sscanf(member, "%d:%d", &userID, &portfolioID)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("failed to parse member: %v", err)
	}

	rebalancingTime := time.Unix(int64(score), 0)

	return userID, portfolioID, rebalancingTime, nil
}
