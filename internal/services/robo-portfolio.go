package services

import (
	"context"
	"fmt"
	"math"
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
	RebalancePortfolio(portfolioID, userID uint) error
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
			go func(asset *models.PortfolioAsset, assetAmount float64) {
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
			}(asset, amount*asset.Percentage/100)
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
func (s *roboPortfolioServiceImpl) RebalancePortfolio(portfolioID, userID uint) error {
	fmt.Println("Rebalancing portfolio", portfolioID, "for user", userID)
	portfolio, err := s.repo.GetPortfolio(portfolioID, userID)
	if err != nil {
		return fmt.Errorf("portfolio %d for user %d returns error: %w", portfolioID, userID, err)
	}

	threshold, exists := commons.RebalancingThresholds[*portfolio.RebalanceFreq]
	if !exists {
		return fmt.Errorf("invalid rebalance frequency: %s", *portfolio.RebalanceFreq)
	}
	//get total portfolio value
	latestAssetPricess := make(map[string]float64)

	totalValue, err := s.getPortfolioValue(portfolio, latestAssetPricess)
	if err != nil {
		return fmt.Errorf("failed to get portfolio category values: %w", err)
	}

	overPerformingAssets := []*models.PortfolioAsset{}
	underPerformingAssets := []*models.PortfolioAsset{}
	var cashCategory *models.PortfolioCategory
	for _, category := range portfolio.Category {
		// handle cash last
		if category.Name == "cash" {
			cashCategory = category
			continue
		}
		for _, asset := range category.Assets {
			latestPrice := latestAssetPricess[asset.Symbol]
			curValue := latestPrice * asset.SharesOwned
			targetValue := totalValue * asset.Percentage / 100
			if (math.Abs(curValue-targetValue) <= targetValue*threshold/100) || (asset.SharesOwned == 0) {
				continue
			} else if curValue > targetValue { // overperforming asset, sell shares
				overPerformingAssets = append(overPerformingAssets, asset)
			} else { // underperforming asset, buy shares
				underPerformingAssets = append(underPerformingAssets, asset)
			}
		}

	}
	totalCash := &cashCategory.TotalAmount

	// correct over-allocated assets by selling shares
	if err := s.sellOverPerformingAssets(overPerformingAssets, latestAssetPricess, totalValue, totalCash); err != nil {
		return err
	}

	// correct under-allocated assets by buying shares
	if err := s.buyUnderPerformingAssets(underPerformingAssets, latestAssetPricess, totalValue, totalCash); err != nil {
		return err
	}

	// balance the cash
	targetCash := totalValue * cashCategory.TotalPercentage / 100

	if *totalCash < targetCash {
		// not enouggh cash, send a notification to the user
		fmt.Printf("Warning: Cash is under-allocated. Expected: %.2f, Available: %.2f\n", targetCash, *totalCash)
		cashCategory.TotalAmount = *totalCash
	} else if *totalCash > targetCash {
		// too much cash: add the extra back into the portfolio
		excessCash := *totalCash - targetCash
		fmt.Printf("Excess cash detected: %.2f. Redistributing...\n", excessCash)
		if err := s.addMoneyToPortfolio(portfolio, excessCash); err != nil {
			return err
		}
		*totalCash = targetCash
	}

	return nil
}

func (s *roboPortfolioServiceImpl) getPortfolioValue(portfolio *models.Portfolio, latestAssetPricess map[string]float64) (float64, error) {
	var mu sync.Mutex

	totalValue := 0.0
	for _, category := range portfolio.Category {
		if category.Name == "cash" {
			totalValue += category.TotalAmount
			continue
		}
		categoryValue := 0.0
		var wg sync.WaitGroup
		errCh := make(chan error, len(category.Assets))
		for _, asset := range category.Assets {
			wg.Add(1)
			go func(asset *models.PortfolioAsset) {
				defer wg.Done()

				latestPrice, err := s.genAIService.GetLatestAssetPrice(asset.Symbol)

				if err != nil {
					errCh <- fmt.Errorf("failed to get latest price for asset %s: %w", asset.Symbol, err)
					return
				}

				assetValue := asset.SharesOwned * latestPrice
				mu.Lock()
				latestAssetPricess[asset.Symbol] = latestPrice
				categoryValue += assetValue
				mu.Unlock()
			}(asset)
		}

		wg.Wait()
		close(errCh)
		for err := range errCh {
			return 0, err
		}
		category.TotalAmount = categoryValue
		totalValue += categoryValue
	}

	return totalValue, nil
}

func (s *roboPortfolioServiceImpl) sellOverPerformingAssets(portfolioCategories []*models.PortfolioAsset, latestAssetPricess map[string]float64, totalValue float64, totalCash *float64) error {
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(portfolioCategories))

	for _, asset := range portfolioCategories {
		wg.Add(1)
		go func(asset *models.PortfolioAsset) {
			defer wg.Done()

			latestPrice, exists := latestAssetPricess[asset.Symbol]
			if !exists {
				errCh <- fmt.Errorf("latest price for asset %s not found", asset.Symbol)
				return
			}

			curValue := latestPrice * asset.SharesOwned
			targetValue := totalValue * asset.Percentage / 100
			amountToSell := curValue - targetValue
			mu.Lock()
			*totalCash += amountToSell
			mu.Unlock()

			sharesToSell := amountToSell / latestPrice
			asset.SharesOwned -= sharesToSell
			asset.TotalInvested -= amountToSell
			asset.AvgBuyPrice = asset.TotalInvested / asset.SharesOwned
		}(asset)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		return err
	}
	return nil
}

func (s *roboPortfolioServiceImpl) buyUnderPerformingAssets(portfolioCategories []*models.PortfolioAsset, latestAssetPricess map[string]float64, totalValue float64, totalCash *float64) error {
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(portfolioCategories))

	for _, asset := range portfolioCategories {
		wg.Add(1)
		go func(asset *models.PortfolioAsset) {
			defer wg.Done()

			latestPrice, exists := latestAssetPricess[asset.Symbol]
			if !exists {
				errCh <- fmt.Errorf("latest price for asset %s not found", asset.Symbol)
				return
			}

			curValue := latestPrice * asset.SharesOwned
			targetValue := totalValue * asset.Percentage / 100
			amountToBuy := targetValue - curValue

			mu.Lock()
			if amountToBuy > *totalCash {
				amountToBuy = *totalCash
			}
			*totalCash -= amountToBuy
			mu.Unlock()

			if amountToBuy == 0 {
				return
			}
			sharesToBuy := amountToBuy / latestPrice
			asset.SharesOwned += sharesToBuy
			asset.TotalInvested += amountToBuy
			asset.AvgBuyPrice = asset.TotalInvested / asset.SharesOwned

		}(asset)
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		return err
	}
	return nil
}
