package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/redis"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

type PortfolioScheduler interface {
	Start(ctx context.Context)
}

type portfolioSchedulerImpl struct {
	ticker  *time.Ticker
	service services.RoboPortfolioService
	repo    repositories.RoboPortfolioRepo
	redis   redis.RoboPortfolioRedis
}

func NewPortfolioSchedulerImpl(s services.RoboPortfolioService, r repositories.RoboPortfolioRepo, c redis.RoboPortfolioRedis) *portfolioSchedulerImpl {
	return &portfolioSchedulerImpl{
		ticker:  time.NewTicker(12 * time.Hour),
		service: s,
		repo:    r,
		redis:   c,
	}
}

func (s *portfolioSchedulerImpl) Start(ctx context.Context) {
	log.Println("Rebalancing portfolios at the start")
	s.rebalancePortfolios(ctx)
	go func() {
		for {
			select {
			case t := <-s.ticker.C:
				log.Println("Rebalancing portfolios at", t)
				s.rebalancePortfolios(ctx)
			case <-ctx.Done():
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *portfolioSchedulerImpl) rebalancePortfolios(ctx context.Context) {
	isEmpty, err := s.redis.IsEmpty(ctx)
	if err != nil {
		log.Println("Failed to check if rebalancing queue is empty:", err)
		return
	}
	if isEmpty {
		log.Println("No portfolios to rebalance")
		return
	}
	portfolioIDs, err := s.redis.GetDuePortfolios(ctx)

	if err != nil {
		log.Println("Failed to get due portfolios for rebalancing:", err)
		return
	}

	const maxWorkers = 100
	workerPool := make(chan struct{}, maxWorkers)
	errChan := make(chan error, len(portfolioIDs))
	var wg sync.WaitGroup

	for _, portfolio := range portfolioIDs {
		var userID, portfolioID uint
		_, err := fmt.Sscanf(portfolio, "%d:%d", &userID, &portfolioID)
		if err != nil {
			log.Printf("Failed to parse portfolio %s: %v", portfolio, err)
			continue
		}
		success, err := s.redis.AcquireLock(ctx, userID, portfolioID, 2*time.Minute)
		if !success || err != nil {
			log.Printf("Portfolio %d:%d is already locked, skipping\n", userID, portfolioID)
			continue
		}

		workerPool <- struct{}{}
		wg.Add(1)
		go func(userID, portfolioID uint) {
			defer func() {
				wg.Done()
				<-workerPool
				if err := s.redis.ReleaseLock(ctx, userID, portfolioID); err != nil {
					log.Printf("Failed to release lock for portfolio %d:%d: %v", userID, portfolioID, err)
				}
			}()

			portfolio, err := s.service.RebalancePortfolio(userID, portfolioID)
			if err != nil {
				errChan <- fmt.Errorf("failed to rebalance portfolio %d:%d: %w", userID, portfolioID, err)
				return
			}

			// delete current portfolio from queue and queue for the next rebalance time
			if err := s.redis.DeletePortfolioFromQueue(ctx, userID, portfolioID); err != nil {
				errChan <- fmt.Errorf("failed to remove portfolio %d:%d from rebalancing queue: %w", userID, portfolioID, err)
				return
			}
			// add the new balance time
			nextRebalanceTime, err := commons.GetNextRebalanceTime(*portfolio.RebalanceFreq)
			if err != nil {
				errChan <- fmt.Errorf("failed to get next rebalance time for portfolio %d:%d: %w", userID, portfolioID, err)
				return
			}

			if err = s.redis.AddPortfolioToRebalancingQueue(ctx, userID, portfolio.ID, nextRebalanceTime); err != nil {
				errChan <- fmt.Errorf("failed to add portfolio %d:%d to rebalancing queue: %w", userID, portfolioID, err)
				return
			}
		}(userID, portfolioID)

	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		log.Println(err)
	}

}
