package services

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type PortfolioService interface {
	ConfirmGeneratedRoboPortfolio(req dto.ConfirmPortfolioRequest, userID uint) error
	GetRoboPortfolio(userID uint) (*models.Portfolio, error)
	AddMoneyToRoboPortfolio(userID uint, amount float64) (*models.Portfolio, error)
	UpdateRebalanceFreq(userID uint, freq string) error

	GetManualPortfolios(userID uint) ([]models.Portfolio, error)
}

type portfolioServiceImpl struct {
	repo           repositories.PortfolioRepo
	profileService ProfileService
	genAIService   GenAIService
}

func NewPortfolioService(pr repositories.PortfolioRepo, ps ProfileService, gs GenAIService) *portfolioServiceImpl {
	return &portfolioServiceImpl{repo: pr, profileService: ps, genAIService: gs}
}

func (s *portfolioServiceImpl) ConfirmGeneratedRoboPortfolio(req dto.ConfirmPortfolioRequest, userID uint) error {
	_, err := s.repo.GetRoboPortfolio(userID)
	if err == nil {
		return err
	}

	portfolio := models.Portfolio{
		UserID:        userID,
		Name:          strconv.FormatUint(uint64(userID), 10) + " Robo Advisor Portfolio",
		IsRoboAdvisor: true,
		Category:      []*models.PortfolioCategory{},
		RebalanceFreq: nil,
	}

	// manually append cash category (no assets)
	cashCategory := &models.PortfolioCategory{
		PortfolioUserID: userID,
		PortfolioID:     portfolio.ID,
		Name:            "cash",
		TotalPercentage: req.Portfolio["cash"],
		TotalAmount:     0,
		Assets:          []*models.PortfolioAsset{},
	}
	portfolio.Category = append(portfolio.Category, cashCategory)

	for categoryName, assets := range req.Allocations {
		category := &models.PortfolioCategory{
			PortfolioUserID: userID,
			PortfolioID:     portfolio.ID,
			Name:            categoryName,
			TotalPercentage: req.Portfolio[categoryName],
			TotalAmount:     0,
			Assets:          []*models.PortfolioAsset{},
		}

		for _, asset := range assets.Assets {
			category.Assets = append(category.Assets, &models.PortfolioAsset{
				PortfolioCategoryID: category.ID,
				Symbol:              asset.Symbol,
				Percentage:          asset.Percentage,
				SharesOwned:         0,
				TotalInvested:       0,
				AvgBuyPrice:         0,
			})
		}

		portfolio.Category = append(portfolio.Category, category)
	}
	return s.repo.CreatePortfolio(&portfolio)
}

func (s *portfolioServiceImpl) GetRoboPortfolio(userID uint) (*models.Portfolio, error) {
	return s.repo.GetRoboPortfolio(userID)
}

func (s *portfolioServiceImpl) AddMoneyToRoboPortfolio(userID uint, amount float64) (*models.Portfolio, error) {
	portfolio, err := s.repo.GetRoboPortfolio(userID)
	if err != nil {
		return nil, err
	}
	return s.addMoneyToPortfolio(portfolio, amount)
}

func (s *portfolioServiceImpl) UpdateRebalanceFreq(userID uint, freq string) error {
	return s.repo.UpdateRebalanceFreq(userID, freq)
}

func (s *portfolioServiceImpl) GetManualPortfolios(userID uint) ([]models.Portfolio, error) {
	return s.repo.GetManualPortfolios(userID)
}

// utility functions
func (s *portfolioServiceImpl) addMoneyToPortfolio(portfolio *models.Portfolio, amount float64) (*models.Portfolio, error) {
	var wg sync.WaitGroup
	errCh := make(chan error, len(portfolio.Category))
	for _, category := range portfolio.Category {
		categoryTotal := amount * category.TotalPercentage / 100
		category.TotalAmount += categoryTotal
		if category.Name == "cash" || category.TotalPercentage == 0 {
			continue
		}
		for _, asset := range category.Assets {
			wg.Add(1)
			go s.updateAsset(asset, amount*asset.Percentage/100, &wg, errCh)
		}
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		return nil, err
	}

	// save the updated portfolio
	if err := s.repo.UpdatePortfolio(portfolio); err != nil {
		return nil, err
	}
	return portfolio, nil
}

func (s *portfolioServiceImpl) updateAsset(asset *models.PortfolioAsset, assetAmount float64, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	latestPrice, err := s.genAIService.GetLatestAssetPrice(asset.Symbol)
	if err != nil {
		errCh <- fmt.Errorf("failed to get latest price for asset %s: %w", asset.Symbol, err)
		return
	}
	asset.SharesOwned += assetAmount / latestPrice
	asset.TotalInvested += assetAmount
	if asset.SharesOwned > 0 {
		asset.AvgBuyPrice = asset.TotalInvested / asset.SharesOwned
	}
}
