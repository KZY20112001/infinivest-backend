package caches

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type PortfolioCache interface {
	AddPortfolioToRebalancingQueue(ctx context.Context, portfolioID, userID uint, nextRebalance time.Time) error
	GetNextRebalanceTime(ctx context.Context, portfolioID, userID uint) (time.Time, error)
	GetEarliestPortfolio(ctx context.Context) (uint, uint, time.Time, error)
	GetDuePortfolios(ctx context.Context) ([]string, error)
	IsEmpty(ctx context.Context) (bool, error)

	AcquireLock(ctx context.Context, portfolioID, userID uint, ttl time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, portfolioID, userID uint) error
}

type portfolioRedis struct {
	client *redis.Client
}

func NewPortfolioRedis(client *redis.Client) *portfolioRedis {
	return &portfolioRedis{client: client}
}

func (r *portfolioRedis) AddPortfolioToRebalancingQueue(ctx context.Context, portfolioID, userID uint, nextRebalance time.Time) error {
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

func (r *portfolioRedis) GetDuePortfolios(ctx context.Context) ([]string, error) {
	key := "rebalancing_queue"
	now := float64(time.Now().Unix())

	results, err := r.client.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", now),
	}).Result()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no portfolios due for rebalancing")
	}

	var portfolios []string
	for _, result := range results {
		member := result.Member.(string) // userID:portfolioID
		portfolios = append(portfolios, member)
	}

	return portfolios, nil
}

func (r *portfolioRedis) IsEmpty(ctx context.Context) (bool, error) {
	key := "rebalancing_queue"
	card, err := r.client.ZCard(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return card == 0, nil
}

func (r *portfolioRedis) AcquireLock(ctx context.Context, portfolioID, userID uint, ttl time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("rebalancing_lock:%d:%d", userID, portfolioID)

	success, err := r.client.SetNX(ctx, lockKey, "locked", ttl).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func (r *portfolioRedis) ReleaseLock(ctx context.Context, portfolioID, userID uint) error {
	lockKey := fmt.Sprintf("rebalancing_lock:%d:%d", userID, portfolioID)
	_, err := r.client.Del(ctx, lockKey).Result()
	return err
}
