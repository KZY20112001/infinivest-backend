package services

import (
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type ManualPortfolioService interface {
	CreateManualPortfolio(userID uint, portfolioName string) error
	GetManualPortfolios(userID uint) ([]*models.ManualPortfolio, error)
	GetManualPortfolio(userID uint, portfolioName string) (*models.ManualPortfolio, error)
}

type manualPortfolioServiceImpl struct {
	repo repositories.ManualPortfolioRepo
}

func NewManualPortfolioService(pr repositories.ManualPortfolioRepo) *manualPortfolioServiceImpl {
	return &manualPortfolioServiceImpl{repo: pr}
}

func (s *manualPortfolioServiceImpl) CreateManualPortfolio(userID uint, portfolioName string) error {
	portfolio := &models.ManualPortfolio{
		UserID:    userID,
		Name:      portfolioName,
		Assets:    []*models.ManualPortfolioAssets{},
		TotalCash: 0,
	}

	return s.repo.CreateManualPortfolio(portfolio)
}

func (s *manualPortfolioServiceImpl) GetManualPortfolios(userID uint) ([]*models.ManualPortfolio, error) {
	return s.repo.GetManualPortfolios(userID)
}

func (s *manualPortfolioServiceImpl) GetManualPortfolio(userID uint, portfolioName string) (*models.ManualPortfolio, error) {
	return s.repo.GetManualPortfolio(userID, portfolioName)
}
