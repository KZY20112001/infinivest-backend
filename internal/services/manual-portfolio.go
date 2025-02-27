package services

import (
	"github.com/KZY20112001/infinivest-backend/internal/caches"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type ManualPortfolioService interface {
	GetManualPortfolios(userID uint) ([]models.RoboPortfolio, error)
}

type manualPortfolioServiceImpl struct {
	repo           repositories.RoboPortfolioRepo
	cache          caches.RoboPortfolioCache
	profileService ProfileService
}

func NewManualPortfolioService(pr repositories.RoboPortfolioRepo, pc caches.RoboPortfolioCache, ps ProfileService) *manualPortfolioServiceImpl {
	return &manualPortfolioServiceImpl{repo: pr, cache: pc, profileService: ps}
}

func (s *manualPortfolioServiceImpl) GetManualPortfolios(userID uint) ([]models.RoboPortfolio, error) {
	return s.repo.GetManualPortfolios(userID)
}
