package scheduler

import (
	"fmt"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/caches"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

type PortfolioScheduler interface {
	Start()
	Stop()
}

type portfolioSchedulerImpl struct {
	ticker  *time.Ticker
	stop    chan struct{}
	service services.RoboPortfolioService
	repo    repositories.PortfolioRepo
	cache   caches.PortfolioCache
}

func NewPortfolioSchedulerImpl(s services.RoboPortfolioService, r repositories.PortfolioRepo, c caches.PortfolioCache) *portfolioSchedulerImpl {
	return &portfolioSchedulerImpl{
		ticker:  time.NewTicker(10 * time.Second),
		stop:    make(chan struct{}),
		service: s,
		repo:    r,
		cache:   c,
	}
}

func (s *portfolioSchedulerImpl) Start() {
	go func() {
		for {
			select {
			case t := <-s.ticker.C:
				s.task(t)
			case <-s.stop:
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *portfolioSchedulerImpl) Stop() {
	close(s.stop)
}

func (s *portfolioSchedulerImpl) task(t time.Time) {
	fmt.Println("CRONJOB: Task executed at:", t)
}
