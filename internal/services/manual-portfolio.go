package services

import (
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type ManualPortfolioService interface {
	CreateManualPortfolio(userID uint, portfolioName string) error
	GetManualPortfolios(userID uint) ([]*models.ManualPortfolio, error)
	GetManualPortfolio(userID uint, portfolioName string) (*models.ManualPortfolio, error)
	AddMoneyToManualPortfolio(userID uint, portfolioName string, amount float64) error
	WithdrawMoneyFromManualPortfolio(userID uint, portfolioName string, amount float64) (float64, error)
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
		Assets:    []*models.ManualPortfolioAsset{},
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

func (s *manualPortfolioServiceImpl) AddMoneyToManualPortfolio(userID uint, portfolioName string, amount float64) error {
	portfolio, err := s.repo.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		return err
	}

	portfolio.TotalCash += amount
	return s.repo.UpdateManualPortfolio(portfolio)
}

func (s *manualPortfolioServiceImpl) WithdrawMoneyFromManualPortfolio(userID uint, portfolioName string, amount float64) (float64, error) {
	portfolio, err := s.repo.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		return 0, nil
	}
	originalAmount := amount
	if portfolio.TotalCash < amount {
		amount -= portfolio.TotalCash
		portfolio.TotalCash = 0
	} else {
		portfolio.TotalCash -= amount
		amount = 0
	}
	if err := s.repo.UpdateManualPortfolio(portfolio); err != nil {
		return 0, err
	}
	return originalAmount - amount, nil
}
