package services

import (
	"github.com/KZY20112001/infinivest-backend/internal/caches"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type ManualPortfolioService interface {
	GetManualPortfolios(userID uint) ([]models.Portfolio, error)
}

type manualPortfolioServiceImpl struct {
	repo           repositories.PortfolioRepo
	cache          caches.PortfolioCache
	profileService ProfileService
}

func NewManualPortfolioService(pr repositories.PortfolioRepo, pc caches.PortfolioCache, ps ProfileService) *manualPortfolioServiceImpl {
	return &manualPortfolioServiceImpl{repo: pr, cache: pc, profileService: ps}
}

func (s *manualPortfolioServiceImpl) GetManualPortfolios(userID uint) ([]models.Portfolio, error) {
	return s.repo.GetManualPortfolios(userID)
}
