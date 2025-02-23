package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/caches"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

type PortfolioScheduler interface {
	Start(ctx context.Context)
}

type portfolioSchedulerImpl struct {
	ticker  *time.Ticker
	service services.RoboPortfolioService
	repo    repositories.PortfolioRepo
	cache   caches.PortfolioCache
}

func NewPortfolioSchedulerImpl(s services.RoboPortfolioService, r repositories.PortfolioRepo, c caches.PortfolioCache) *portfolioSchedulerImpl {
	return &portfolioSchedulerImpl{
		ticker:  time.NewTicker(1 * time.Minute),
		service: s,
		repo:    r,
		cache:   c,
	}
}

func (s *portfolioSchedulerImpl) Start(ctx context.Context) {
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
	isEmpty, err := s.cache.IsEmpty(ctx)
	if err != nil {
		log.Println("Failed to check if rebalancing queue is empty:", err)
		return
	}
	if isEmpty {
		log.Println("No portfolios to rebalance")
		return
	}
	portfolioIDs, err := s.cache.GetDuePortfolios(ctx)

	if err != nil {
		log.Println("Failed to get due portfolios for rebalancing:", err)
		return
	}

	const maxWorkers = 100
	workerPool := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, portfolio := range portfolioIDs {
		var userID, portfolioID uint
		_, err := fmt.Sscanf(portfolio, "%d:%d", &userID, &portfolioID)
		if err != nil {
			log.Printf("Failed to parse portfolio %s: %v", portfolio, err)
			continue
		}
		success, err := s.cache.AcquireLock(ctx, userID, portfolioID, 2*time.Minute)
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
				if err := s.cache.ReleaseLock(ctx, userID, portfolioID); err != nil {
					log.Printf("Failed to release lock for portfolio %d:%d: %v", userID, portfolioID, err)
				}
			}()

			s.service.RebalancePortfolio(userID, portfolioID)
		}(userID, portfolioID)

	}
}
