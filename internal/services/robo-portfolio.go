package services

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/KZY20112001/infinivest-backend/internal/caches"
	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type RoboPortfolioService interface {
	ConfirmGeneratedRoboPortfolio(req dto.ConfirmPortfolioRequest, userID uint) error
	GetRoboPortfolio(userID uint) (*models.Portfolio, error)
	AddMoneyToRoboPortfolio(ctx context.Context, userID uint, amount float64) (*models.Portfolio, error)
	UpdateRebalanceFreq(userID uint, freq string) error
}

type roboPortfolioServiceImpl struct {
	repo           repositories.PortfolioRepo
	cache          caches.PortfolioCache
	profileService ProfileService
	genAIService   GenAIService
}

func NewRoboPortfolioService(pr repositories.PortfolioRepo, pc caches.PortfolioCache, ps ProfileService, gs GenAIService) *roboPortfolioServiceImpl {
	return &roboPortfolioServiceImpl{repo: pr, cache: pc, profileService: ps, genAIService: gs}
}

func (s *roboPortfolioServiceImpl) ConfirmGeneratedRoboPortfolio(req dto.ConfirmPortfolioRequest, userID uint) error {
	_, err := s.repo.GetRoboPortfolio(userID)
	if err == nil {
		return err
	}

	portfolio := models.Portfolio{
		UserID:        userID,
		Name:          strconv.FormatUint(uint64(userID), 10) + " Robo Advisor Portfolio",
		IsRoboAdvisor: true,
		Category:      []*models.PortfolioCategory{},
		RebalanceFreq: &req.Frequency,
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

func (s *roboPortfolioServiceImpl) GetRoboPortfolio(userID uint) (*models.Portfolio, error) {
	return s.repo.GetRoboPortfolio(userID)
}

func (s *roboPortfolioServiceImpl) AddMoneyToRoboPortfolio(ctx context.Context, userID uint, amount float64) (*models.Portfolio, error) {
	portfolio, err := s.repo.GetRoboPortfolio(userID)
	if err != nil {
		return nil, err
	}

	if err := s.addMoneyToPortfolio(portfolio, amount); err != nil {
		return nil, err
	}

	// check if current portfolio is in queue
	if _, err := s.cache.GetNextRebalanceTime(ctx, portfolio.ID, userID); err == nil {
		return portfolio, nil
	}

	// add portfolio to rebalancing queue
	nextRebalanceTime, err := commons.GetNextRebalanceTime(*portfolio.RebalanceFreq)
	if err != nil {
		return nil, err
	}

	if err = s.cache.AddPortfolioToRebalancingQueue(ctx, portfolio.ID, userID, nextRebalanceTime); err != nil {
		return nil, err
	}
	return portfolio, err
}

func (s *roboPortfolioServiceImpl) UpdateRebalanceFreq(userID uint, freq string) error {
	return s.repo.UpdateRebalanceFreq(userID, freq)
}

// utility functions
func (s *roboPortfolioServiceImpl) addMoneyToPortfolio(portfolio *models.Portfolio, amount float64) error {
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
		return err
	}

	// save the updated portfolio
	if err := s.repo.UpdatePortfolio(portfolio); err != nil {
		return err
	}
	return nil
}

func (s *roboPortfolioServiceImpl) updateAsset(asset *models.PortfolioAsset, assetAmount float64, wg *sync.WaitGroup, errCh chan<- error) {
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
