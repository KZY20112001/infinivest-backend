package services

import (
	"fmt"
	"sync"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type ManualPortfolioService interface {
	GetManualPortfoliosDetails(userID uint) ([]*models.ManualPortfolio, error)
	GetManualPortfoliosSummaries(userID uint) ([]dto.ManualPortfolioSummaryResponse, error)

	GetManualPortfolio(userID uint, portfolioName string) (*models.ManualPortfolio, error)
	GetPorfolioValue(userID uint, portfolioName string) (float64, error)

	CreateManualPortfolio(userID uint, portfolioName string) error
	UpdatePortfolioName(userID uint, portfolioName, newName string) error

	AddMoneyToManualPortfolio(userID uint, portfolioName string, amount float64) error
	WithdrawMoneyFromManualPortfolio(userID uint, portfolioName string, amount float64) (float64, error)

	BuyAssetForManualPortfolio(userID uint, portfolioName, symbol string, shares float64) error
	SellAssetForManualPortfolio(userID uint, portfolioName, symbol string, shares float64) error

	DeleteManualPortfolio(userID uint, portfolioName string) error
}

type manualPortfolioServiceImpl struct {
	repo         repositories.ManualPortfolioRepo
	genAiService GenAIService
}

func NewManualPortfolioService(pr repositories.ManualPortfolioRepo, gs GenAIService) *manualPortfolioServiceImpl {
	return &manualPortfolioServiceImpl{repo: pr, genAiService: gs}
}

func (s *manualPortfolioServiceImpl) GetManualPortfoliosDetails(userID uint) ([]*models.ManualPortfolio, error) {
	return s.repo.GetManualPortfolios(userID)
}

func (s *manualPortfolioServiceImpl) GetManualPortfolio(userID uint, portfolioName string) (*models.ManualPortfolio, error) {
	return s.repo.GetManualPortfolio(userID, portfolioName)
}

func (s *manualPortfolioServiceImpl) GetManualPortfoliosSummaries(userID uint) ([]dto.ManualPortfolioSummaryResponse, error) {
	portfolios, err := s.repo.GetManualPortfolios(userID)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	summaries := make(chan dto.ManualPortfolioSummaryResponse, len(portfolios))
	for _, portfolio := range portfolios {
		wg.Add(1)
		go func(portfolio *models.ManualPortfolio) {
			defer wg.Done()
			totalValue, err := s.GetPorfolioValue(portfolio.UserID, portfolio.Name)
			if err != nil {
				return
			}
			summaries <- dto.ManualPortfolioSummaryResponse{
				Name:       portfolio.Name,
				TotalValue: totalValue,
			}
		}(portfolio)
	}
	wg.Wait()
	close(summaries)
	var res []dto.ManualPortfolioSummaryResponse
	for summary := range summaries {
		res = append(res, summary)
	}
	return res, nil
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

func (s *manualPortfolioServiceImpl) GetPorfolioValue(userID uint, portfolioName string) (float64, error) {
	portfolio, err := s.repo.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		return 0, err
	}

	totalValue := portfolio.TotalCash
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, asset := range portfolio.Assets {
		wg.Add(1)
		go func(asset *models.ManualPortfolioAsset) {
			defer wg.Done()
			latestValue, err := s.genAiService.GetLatestAssetPrice(asset.Symbol)
			if err != nil {
				return
			}
			mu.Lock()
			totalValue += asset.SharesOwned * latestValue
			mu.Unlock()
		}(asset)

	}
	wg.Wait()
	return totalValue, nil
}

func (s *manualPortfolioServiceImpl) AddMoneyToManualPortfolio(userID uint, portfolioName string, amount float64) error {
	portfolio, err := s.repo.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		return err
	}

	portfolio.TotalCash += amount
	return s.repo.UpdateManualPortfolio(portfolio)
}

func (s *manualPortfolioServiceImpl) UpdatePortfolioName(userID uint, portfolioName, newName string) error {
	portfolio, err := s.repo.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		return err
	}
	return s.repo.UpdateManualPortfolioName(portfolio, newName)
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

func (s *manualPortfolioServiceImpl) BuyAssetForManualPortfolio(userID uint, portfolioName, symbol string, shares float64) error {
	portfolio, err := s.repo.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		return err
	}

	latestValue, err := s.genAiService.GetLatestAssetPrice(symbol)
	if err != nil {
		return err
	}

	totalCost := latestValue * shares
	if totalCost > portfolio.TotalCash {
		return fmt.Errorf("insufficient funds to buy %f shares of %s", shares, symbol)
	}

	found := false
	var curAsset *models.ManualPortfolioAsset
	for _, asset := range portfolio.Assets {
		if asset.Symbol == symbol {
			found = true
			curAsset = asset
			break
		}
	}
	if !found {
		curAsset = &models.ManualPortfolioAsset{
			Symbol:        symbol,
			SharesOwned:   shares,
			TotalInvested: totalCost,
			AvgBuyPrice:   latestValue,
		}
		portfolio.Assets = append(portfolio.Assets, curAsset)
	} else {
		curAsset.SharesOwned += shares
		curAsset.TotalInvested += totalCost
		curAsset.AvgBuyPrice = curAsset.TotalInvested / curAsset.SharesOwned
	}
	portfolio.TotalCash -= totalCost

	return s.repo.UpdateManualPortfolio(portfolio)
}

func (s *manualPortfolioServiceImpl) SellAssetForManualPortfolio(userID uint, portfolioName, symbol string, shares float64) error {
	portfolio, err := s.repo.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		return err
	}
	found := false
	var curAsset *models.ManualPortfolioAsset
	for _, asset := range portfolio.Assets {
		if asset.Symbol == symbol {
			found = true
			curAsset = asset
			break
		}
	}
	if !found {
		return fmt.Errorf("asset not found in portfolio")
	}

	if curAsset.SharesOwned < shares {
		return fmt.Errorf("insufficient shares to sell")
	}
	latestValue, err := s.genAiService.GetLatestAssetPrice(symbol)
	if err != nil {
		return err
	}
	portfolio.TotalCash += latestValue * shares
	curAsset.SharesOwned -= shares
	return s.repo.UpdateManualPortfolio(portfolio)
}

func (s *manualPortfolioServiceImpl) DeleteManualPortfolio(userID uint, portfolioName string) error {
	portfolio, err := s.repo.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		return err
	}
	totalValue, err := s.GetPorfolioValue(userID, portfolioName)
	if err != nil {
		return err
	}
	if totalValue > 0 {
		return fmt.Errorf("%s portfolio has assets and liquid cash. Sell all assets and withdraw the money before deleting the portfolio", portfolioName)
	}
	return s.repo.DeleteManualPortfolio(portfolio)
}
