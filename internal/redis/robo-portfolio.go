package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RoboPortfolioRedis interface {
	AddPortfolioToRebalancingQueue(ctx context.Context, userID, portfolioID uint, nextRebalance time.Time) error
	GetNextRebalanceTime(ctx context.Context, userID, portfolioID uint) (time.Time, error)
	DeletePortfolioFromQueue(ctx context.Context, userID, portfolioID uint) error
	GetDuePortfolios(ctx context.Context) ([]string, error)
	IsEmpty(ctx context.Context) (bool, error)

	AcquireLock(ctx context.Context, userID, portfolioID uint, ttl time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, userID, portfolioID uint) error

	SetLastSeen(ctx context.Context, userID, portfolioID uint) error
	GetLastSeen(ctx context.Context, userID, portfolioID uint) (time.Time, error)
}

type roboPortfolioRedis struct {
	client *redis.Client
}

func NewRoboPortfolioRedis(client *redis.Client) *roboPortfolioRedis {
	return &roboPortfolioRedis{client: client}
}

func (r *roboPortfolioRedis) AddPortfolioToRebalancingQueue(ctx context.Context, userID, portfolioID uint, nextRebalance time.Time) error {
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

func (r *roboPortfolioRedis) GetNextRebalanceTime(ctx context.Context, userID, portfolioID uint) (time.Time, error) {
	key := "rebalancing_queue"
	member := fmt.Sprintf("%d:%d", userID, portfolioID)

	nextRebalanceTime, err := r.client.ZScore(ctx, key, member).Result()

	if err != nil {
		if err == redis.Nil {
			return time.Time{}, fmt.Errorf("portfolio %d for user %d not found in rebalancing queue", userID, portfolioID)
		}
		return time.Time{}, err
	}
	return time.Unix(int64(nextRebalanceTime), 0), nil
}

func (r *roboPortfolioRedis) DeletePortfolioFromQueue(ctx context.Context, userID, portfolioID uint) error {
	key := "rebalancing_queue"
	member := fmt.Sprintf("%d:%d", userID, portfolioID)
	_, err := r.client.ZRem(ctx, key, member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *roboPortfolioRedis) GetDuePortfolios(ctx context.Context) ([]string, error) {

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
	fmt.Println("Due portfolios:", results)
	var portfolios []string
	for _, result := range results {
		member := result.Member.(string) // userID:portfolioID
		portfolios = append(portfolios, member)
	}

	return portfolios, nil
}

func (r *roboPortfolioRedis) IsEmpty(ctx context.Context) (bool, error) {
	key := "rebalancing_queue"
	card, err := r.client.ZCard(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return card == 0, nil
}

func (r *roboPortfolioRedis) SetLastSeen(ctx context.Context, userID, portfolioID uint) error {
	key := fmt.Sprintf("rebalance_last_seen:%d:%d", userID, portfolioID)
	_, err := r.client.Set(ctx, key, time.Now().Unix(), 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *roboPortfolioRedis) GetLastSeen(ctx context.Context, userID, portfolioID uint) (time.Time, error) {
	key := fmt.Sprintf("rebalance_last_seen:%d:%d", userID, portfolioID)
	lastSeen, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil { // zero time returned if key does not exist yet
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	lastSeenUnix, err := strconv.ParseInt(lastSeen, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(lastSeenUnix, 0), nil
}

func (r *roboPortfolioRedis) AcquireLock(ctx context.Context, userID, portfolioID uint, ttl time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("rebalancing_lock:%d:%d", userID, portfolioID)

	success, err := r.client.SetNX(ctx, lockKey, "locked", ttl).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func (r *roboPortfolioRedis) ReleaseLock(ctx context.Context, userID, portfolioID uint) error {
	lockKey := fmt.Sprintf("rebalancing_lock:%d:%d", userID, portfolioID)
	_, err := r.client.Del(ctx, lockKey).Result()
	return err
}
