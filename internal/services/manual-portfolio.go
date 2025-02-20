package services

import (
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type ManualPortfolioService interface {
	GetManualPortfolios(userID uint) ([]models.Portfolio, error)
}

type manualPortfolioServiceImpl struct {
	repo repositories.PortfolioRepo
}

func NewManualPortfolioService(pr repositories.PortfolioRepo) *manualPortfolioServiceImpl {
	return &manualPortfolioServiceImpl{repo: pr}
}

func (s *manualPortfolioServiceImpl) GetManualPortfolios(userID uint) ([]models.Portfolio, error) {
	return s.repo.GetManualPortfolios(userID)
}
