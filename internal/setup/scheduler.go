package setup

import (
	"github.com/KZY20112001/infinivest-backend/internal/redis"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/KZY20112001/infinivest-backend/internal/scheduler"
	"github.com/KZY20112001/infinivest-backend/internal/services"
)

func PortfolioScheduler(
	roboPortfolioService services.RoboPortfolioService,
	portfolioRepo repositories.RoboPortfolioRepo,
	portfolioRedis redis.RoboPortfolioRedis) scheduler.PortfolioScheduler {
	return scheduler.NewPortfolioSchedulerImpl(
		roboPortfolioService,
		portfolioRepo,
		portfolioRedis,
	)
}
