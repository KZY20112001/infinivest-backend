package setup

import (
	"github.com/KZY20112001/infinivest-backend/internal/caches"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/scheduler"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

func PortfolioScheduler(
	roboPortfolioService services.RoboPortfolioService,
	portfolioRepo repositories.RoboPortfolioRepo,
	portfolioCache caches.RoboPortfolioCache) scheduler.PortfolioScheduler {
	return scheduler.NewPortfolioSchedulerImpl(
		roboPortfolioService,
		portfolioRepo,
		portfolioCache,
	)
}
