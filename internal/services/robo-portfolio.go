package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/commons/email"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/redis"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type RoboPortfolioService interface {
	ConfirmGeneratedRoboPortfolio(req dto.ConfirmPortfolioRequest, userID uint) error
	GetRoboPortfolioDetails(userID uint) (*models.RoboPortfolio, error)
	GetRoboPortfolioSummary(userID uint) (dto.RoboPortfolioSummaryResponse, error)
	AddMoneyToRoboPortfolio(ctx context.Context, userID uint, amount float64) (*models.RoboPortfolio, error)
	WithDrawMoneyFromRoboPortfolio(ctx context.Context, userID uint, amount float64) (float64, error)
	UpdateRebalanceFreq(ctx context.Context, userID uint, freq string) error
	RebalancePortfolio(ctx context.Context, userID, portfolioID uint) (*models.RoboPortfolio, error)
	DeleteRoboPortfolio(ctx context.Context, userID uint) error

	GetRoboPortfolioTransactions(userID uint, limit int) ([]*models.RoboPortfolioTransaction, error)

	GetRebalanceEvents(ctx context.Context, userID uint) ([]*models.RebalanceEvent, error)
	UpdateLastSeenRebalanceTime(ctx context.Context, userID uint) error
}

type roboPortfolioServiceImpl struct {
	repo                repositories.RoboPortfolioRepo
	redis               redis.RoboPortfolioRedis
	genAIService        GenAIService
	notificationService NotificationService
	userService         UserService
}

func NewRoboPortfolioService(pr repositories.RoboPortfolioRepo, pc redis.RoboPortfolioRedis, gs GenAIService, ns NotificationService, us UserService) *roboPortfolioServiceImpl {
	return &roboPortfolioServiceImpl{repo: pr, redis: pc, genAIService: gs, notificationService: ns, userService: us}
}

func (s *roboPortfolioServiceImpl) ConfirmGeneratedRoboPortfolio(req dto.ConfirmPortfolioRequest, userID uint) error {
	_, err := s.repo.GetRoboPortfolioDetails(userID)
	if err == nil {
		return fmt.Errorf("roboportfolio already exists for user %d", userID)
	}

	portfolio := models.RoboPortfolio{
		UserID:        userID,
		Categories:    []*models.RoboPortfolioCategory{},
		RebalanceFreq: &req.Frequency,
		IsRebalancing: false,
	}

	// manually append cash category (no assets)
	cashCategory := &models.RoboPortfolioCategory{
		RoboPortfolioID: portfolio.ID,
		Name:            "cash",
		TotalPercentage: req.Portfolio["cash"],
		TotalAmount:     0,
		Assets:          []*models.RoboPortfolioAsset{},
	}
	portfolio.Categories = append(portfolio.Categories, cashCategory)

	for categoryName, assets := range req.Allocations {
		category := &models.RoboPortfolioCategory{
			RoboPortfolioID: portfolio.ID,
			Name:            categoryName,
			TotalPercentage: req.Portfolio[categoryName],
			TotalAmount:     0,
			Assets:          []*models.RoboPortfolioAsset{},
		}

		for _, asset := range assets.Assets {
			category.Assets = append(category.Assets, &models.RoboPortfolioAsset{
				RoboPortfolioCategoryID: category.ID,
				Name:                    asset.Name,
				Symbol:                  asset.Symbol,
				Percentage:              asset.Percentage,
				SharesOwned:             0,
				TotalInvested:           0,
				AvgBuyPrice:             0,
			})
		}

		portfolio.Categories = append(portfolio.Categories, category)
	}
	return s.repo.CreateRoboPortfolio(&portfolio)
}

func (s *roboPortfolioServiceImpl) GetRoboPortfolioDetails(userID uint) (*models.RoboPortfolio, error) {
	return s.repo.GetRoboPortfolioDetails(userID)
}

func (s *roboPortfolioServiceImpl) GetRoboPortfolioSummary(userID uint) (dto.RoboPortfolioSummaryResponse, error) {
	portfolio, err := s.repo.GetRoboPortfolioDetails(userID)
	if err != nil {
		return dto.RoboPortfolioSummaryResponse{}, err
	}
	totalValue, err := s.getPortfolioValue(portfolio, make(map[string]float64))
	if err != nil {
		return dto.RoboPortfolioSummaryResponse{}, err
	}
	return dto.RoboPortfolioSummaryResponse{
		RebalanceFreq: *portfolio.RebalanceFreq,
		TotalValue:    totalValue}, nil
}

func (s *roboPortfolioServiceImpl) AddMoneyToRoboPortfolio(ctx context.Context, userID uint, amount float64) (*models.RoboPortfolio, error) {
	portfolio, err := s.repo.GetRoboPortfolioDetails(userID)
	if err != nil {
		return nil, err
	}

	// lock the portfolio to prevent concurrent updates
	if err := s.repo.LockRoboPortfolio(portfolio); err != nil {
		return nil, err
	}

	if err := s.addMoneyToPortfolio(portfolio, amount); err != nil {
		return nil, err
	}

	err = s.repo.CreateRoboPortfolioTransaction(&models.RoboPortfolioTransaction{
		RoboPortfolioID: portfolio.ID,
		TransactionType: "deposit",
		TotalAmount:     amount,
	})

	if err != nil {
		return nil, err
	}
	if err := s.repo.UnlockRoboPortfolio(portfolio); err != nil {
		return nil, err
	}

	// check if current portfolio is in queue
	if _, err := s.redis.GetNextRebalanceTime(ctx, userID, portfolio.ID); err == nil {
		return portfolio, nil
	}

	// add portfolio to rebalancing queue
	nextRebalanceTime, err := commons.GetNextRebalanceTime(*portfolio.RebalanceFreq)
	if err != nil {
		return nil, err
	}

	if err = s.redis.AddPortfolioToRebalancingQueue(ctx, userID, portfolio.ID, nextRebalanceTime); err != nil {
		return nil, err
	}

	return portfolio, nil
}

func (s *roboPortfolioServiceImpl) WithDrawMoneyFromRoboPortfolio(ctx context.Context, userID uint, amount float64) (float64, error) {
	portfolio, err := s.repo.GetRoboPortfolioDetails(userID)
	if err != nil {
		return 0, err
	}

	// lock the portfolio to prevent concurrent updates
	if err := s.repo.LockRoboPortfolio(portfolio); err != nil {
		return 0, err
	}

	var cashCategory *models.RoboPortfolioCategory

	for _, category := range portfolio.Categories {
		if category.Name == "cash" {
			cashCategory = category
			break
		}
	}
	if amount <= cashCategory.TotalAmount {
		cashCategory.TotalAmount -= amount
		if err := s.repo.UpdateRoboPortfolio(portfolio); err != nil {
			return 0, err
		}
		return amount, nil
	}

	// sell assets to cover the remaining amount
	leftToWithdraw := amount
	withdrawn := 0.0
	withdrawn += cashCategory.TotalAmount
	amount -= cashCategory.TotalAmount
	cashCategory.TotalAmount = 0
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, category := range portfolio.Categories {
		if category.Name == "cash" || category.TotalPercentage == 0 {
			continue
		}

		errChan := make(chan error, len(category.Assets))
		for _, asset := range category.Assets {
			wg.Add(1)
			go func(asset *models.RoboPortfolioAsset) {
				defer wg.Done()
				latestPrice, err := s.genAIService.GetLatestAssetPrice(asset.Symbol)
				if err != nil {
					errChan <- fmt.Errorf("failed to get latest price for asset %s: %w", asset.Symbol, err)
					return
				}
				amountToSell := amount * asset.Percentage / 100
				numOfShares := amountToSell / latestPrice
				if numOfShares > asset.SharesOwned {
					numOfShares = asset.SharesOwned
				}
				curAmount := numOfShares * latestPrice
				mu.Lock()
				category.TotalAmount -= curAmount
				if category.TotalAmount < 0 {
					category.TotalAmount = 0
				}
				leftToWithdraw -= curAmount
				withdrawn += curAmount
				mu.Unlock()
				asset.SharesOwned -= numOfShares
				asset.TotalInvested -= curAmount
				if asset.TotalInvested < 0 {
					asset.TotalInvested = 0
				}
				if asset.SharesOwned > 0 {
					asset.AvgBuyPrice = asset.TotalInvested / asset.SharesOwned
				}
			}(asset)
		}
		wg.Wait()
		close(errChan)
		for err := range errChan {
			log.Println(err)
		}
	}

	if err := s.repo.UpdateRoboPortfolio(portfolio); err != nil {
		return 0, err
	}

	// unlock the portfolio after adding to queue
	if err := s.repo.UnlockRoboPortfolio(portfolio); err != nil {
		return 0, err
	}

	err = s.repo.CreateRoboPortfolioTransaction(&models.RoboPortfolioTransaction{
		RoboPortfolioID: portfolio.ID,
		TransactionType: "withdrawal",
		TotalAmount:     withdrawn,
	})

	if err != nil {
		return 0, err
	}

	return withdrawn, nil
}

func (s *roboPortfolioServiceImpl) UpdateRebalanceFreq(ctx context.Context, userID uint, freq string) error {
	portfolio, err := s.repo.GetRoboPortfolioDetails(userID)
	if err != nil {
		return err
	}

	if portfolio.RebalanceFreq != nil && *portfolio.RebalanceFreq == freq {
		return nil
	}
	if _, err := s.RebalancePortfolio(ctx, userID, portfolio.ID); err != nil {
		return err
	}

	if err := s.redis.DeletePortfolioFromQueue(ctx, userID, portfolio.ID); err != nil {
		return err
	}

	if err := s.repo.UpdateRebalanceFreq(userID, freq); err != nil {
		return err
	}

	nextRebalanceTime, err := commons.GetNextRebalanceTime(freq)
	if err != nil {
		return err
	}
	return s.redis.AddPortfolioToRebalancingQueue(ctx, userID, portfolio.ID, nextRebalanceTime)
}

func (s *roboPortfolioServiceImpl) DeleteRoboPortfolio(ctx context.Context, userID uint) error {
	portfolio, err := s.repo.GetRoboPortfolioDetails(userID)
	if err != nil {
		return err
	}

	totalValue, err := s.getPortfolioValue(portfolio, make(map[string]float64))
	if err != nil {
		return err
	}
	if totalValue > 0 {
		return fmt.Errorf("portfolio %d for user %d has a total value of %.2f. Please withdraw all funds before deleting", portfolio.ID, userID, totalValue)
	}
	if err := s.redis.DeletePortfolioFromQueue(ctx, userID, portfolio.ID); err != nil {
		return err
	}
	return s.repo.DeleteRoboPortfolio(portfolio)
}

func (s *roboPortfolioServiceImpl) RebalancePortfolio(ctx context.Context, userID, portfolioID uint) (*models.RoboPortfolio, error) {
	log.Println("Rebalancing portfolio", portfolioID, "for user", userID)
	portfolio, err := s.repo.GetRoboPortfolioDetails(userID)
	if err != nil {
		return nil, fmt.Errorf("portfolio %d for user %d returns error: %w", userID, portfolioID, err)
	}
	// lock the portfolio to prevent concurrent updates
	if err := s.repo.LockRoboPortfolio(portfolio); err != nil {
		return nil, err
	}
	threshold, exists := commons.RebalancingThresholds[*portfolio.RebalanceFreq]
	if !exists {
		return nil, fmt.Errorf("invalid rebalance frequency: %s", *portfolio.RebalanceFreq)
	}

	//get total portfolio value
	latestAssetPrices := make(map[string]float64)
	totalValue, err := s.getPortfolioValue(portfolio, latestAssetPrices)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio value: %w", err)
	}

	overPerformingAssets := []*models.RoboPortfolioAsset{}
	underPerformingAssets := []*models.RoboPortfolioAsset{}
	var cashCategory *models.RoboPortfolioCategory
	for _, category := range portfolio.Categories {
		// handle cash last
		if category.Name == "cash" {
			cashCategory = category
			continue
		}
		for _, asset := range category.Assets {
			latestPrice := latestAssetPrices[asset.Symbol]
			curValue := latestPrice * asset.SharesOwned
			targetValue := totalValue * asset.Percentage / 100
			if (math.Abs(curValue-targetValue) <= targetValue*threshold/100) || (asset.SharesOwned == 0) {
				log.Println("Asset", asset.Symbol, "is within threshold")
				continue
			} else if curValue > targetValue { // overperforming asset, sell shares
				overPerformingAssets = append(overPerformingAssets, asset)
			} else { // underperforming asset, buy shares
				underPerformingAssets = append(underPerformingAssets, asset)
			}
		}

	}

	totalCash := &cashCategory.TotalAmount
	totalSellAmount, err := s.sellOverPerformingAssets(overPerformingAssets, latestAssetPrices, totalValue, totalCash)
	if err != nil {
		return nil, err
	}

	var failReason string = ""
	totalBuyAmount, err := s.buyUnderPerformingAssets(underPerformingAssets, latestAssetPrices, totalValue, totalCash, &failReason)
	if err != nil {
		return nil, err
	}

	// balance the cash
	targetCash := totalValue * cashCategory.TotalPercentage / 100

	if *totalCash < targetCash {
		cashCategory.TotalAmount = *totalCash

		// not enough cash, send a notification to the user
		failReason = "Cash is under-allocated. Expected: " + fmt.Sprintf("%.2f", targetCash) + ", Available: " + fmt.Sprintf("%.2f", *totalCash)
		if err := s.notificationService.AddNotification(ctx, userID, "rebalance", failReason); err != nil {
			log.Println("Failed to add notification:", err)
		}
		log.Printf("Warning: Cash is under-allocated. Expected: %.2f, Available: %.2f\n", targetCash, *totalCash)
		user, err := s.userService.GetUser(userID)
		if err != nil {
			log.Printf("Failed to get profile for user %d: %v\n", userID, err)
		} else {
			subject := "Alert from InfiniVest"
			body := fmt.Sprintf(`
			<p>Dear User,</p>
			<p>We encountered an issue during the rebalancing process of your robo-portfolio due to insufficient available cash.</p>
			<p><strong>Expected Cash:</strong> $%.2f<br>
			<strong>Available Cash:</strong> $%.2f</p>
			<p>Please top up your robo-portfolio to ensure optimal performance.</p>
			<p>If you have any questions, feel free to reach out to our support team.</p>
			<p>Best regards,<br>
			The InfiniVest Team</p>
		`, targetCash, *totalCash)
			if err := email.SendEmail(user.Email, subject, body); err != nil {
				log.Println("Failed to send email:", err)
			}
		}
	} else if *totalCash > targetCash {
		// too much cash: add the extra back into the portfolio
		excessCash := *totalCash - targetCash
		log.Printf("Excess cash detected: %.2f. Redistributing...\n", excessCash)
		if err := s.addMoneyToPortfolio(portfolio, excessCash); err != nil {
			return nil, err
		}
		*totalCash = targetCash
	}

	var reason *string = nil
	if failReason != "" {
		reason = &failReason
	}
	rebalanceEvent := &models.RebalanceEvent{
		RoboPortfolioID:      portfolio.ID,
		PortfolioValueBefore: totalValue,
		PortfolioValueAfter:  totalValue,
		TotalBuyAmount:       totalBuyAmount,
		TotalSellAmount:      totalSellAmount,
		NetChange:            0,
		Success:              failReason == "",
		Reason:               reason,
	}
	if err := s.repo.UpdateRoboPortfolio(portfolio); err != nil {
		return nil, err
	}

	if err := s.repo.CreateRebalanceEvent(rebalanceEvent); err != nil {
		return nil, err
	}

	log.Println("Rebalanced portfolio:", portfolioID, "for user", userID)
	if err := s.repo.UnlockRoboPortfolio(portfolio); err != nil {
		return nil, err
	}

	// add to notification queue
	var message string = "Portfolio rebalanced successfully!"
	if failReason != "" {
		message = "Portfolio rebalanced with issues!"
		if err := s.notificationService.AddNotification(ctx, userID, "rebalance", failReason); err != nil {
			log.Println("Failed to add notification:", err)
		}
	}

	if err := s.notificationService.AddNotification(ctx, userID, "rebalance", message); err != nil {
		log.Println("Failed to add notification:", err)

	}
	return portfolio, nil
}

func (s *roboPortfolioServiceImpl) GetRoboPortfolioTransactions(userID uint, limit int) ([]*models.RoboPortfolioTransaction, error) {
	return s.repo.GetRoboPortfolioTransactions(userID, limit)
}

func (s *roboPortfolioServiceImpl) GetRebalanceEvents(ctx context.Context, userID uint) ([]*models.RebalanceEvent, error) {
	portfolio, err := s.repo.GetRoboPortfolioDetails(userID)
	if err != nil {
		return nil, fmt.Errorf("portfolio not found for user %d", userID)
	}
	lastSeen, err := s.redis.GetLastSeen(ctx, userID, portfolio.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last seen time for portfolio %d for user %d: %w", portfolio.ID, userID, err)
	}
	return s.repo.GetRebalanceEvents(userID, lastSeen)
}

func (s *roboPortfolioServiceImpl) UpdateLastSeenRebalanceTime(ctx context.Context, userID uint) error {
	portfolio, err := s.repo.GetRoboPortfolioDetails(userID)
	if err != nil {
		return fmt.Errorf("portfolio not found for user %d", userID)
	}
	return s.redis.SetLastSeen(ctx, userID, portfolio.ID)
}

func (s *roboPortfolioServiceImpl) addMoneyToPortfolio(portfolio *models.RoboPortfolio, amount float64) error {
	var wg sync.WaitGroup
	for _, category := range portfolio.Categories {
		categoryTotal := amount * category.TotalPercentage / 100
		category.TotalAmount += categoryTotal
		if category.Name == "cash" || category.TotalPercentage == 0 {
			continue
		}
		errCh := make(chan error, len(category.Assets))
		for _, asset := range category.Assets {
			wg.Add(1)
			go func(asset *models.RoboPortfolioAsset, assetAmount float64) {
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
		wg.Wait()
		close(errCh)
		for err := range errCh {
			return err
		}
	}

	// save the updated portfolio
	if err := s.repo.UpdateRoboPortfolio(portfolio); err != nil {
		return err
	}
	return nil
}

func (s *roboPortfolioServiceImpl) getPortfolioValue(portfolio *models.RoboPortfolio, latestAssetPrices map[string]float64) (float64, error) {
	var mu sync.Mutex

	totalValue := 0.0
	for _, category := range portfolio.Categories {
		if category.Name == "cash" {
			totalValue += category.TotalAmount
			continue
		}
		categoryValue := 0.0
		var wg sync.WaitGroup
		errCh := make(chan error, len(category.Assets))
		for _, asset := range category.Assets {
			wg.Add(1)
			go func(asset *models.RoboPortfolioAsset) {
				defer wg.Done()

				latestPrice, err := s.genAIService.GetLatestAssetPrice(asset.Symbol)
				if err != nil {
					errCh <- fmt.Errorf("failed to get latest price for asset %s: %w", asset.Symbol, err)
					return
				}

				assetValue := asset.SharesOwned * latestPrice
				mu.Lock()
				latestAssetPrices[asset.Symbol] = latestPrice
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

func (s *roboPortfolioServiceImpl) sellOverPerformingAssets(portfolioCategories []*models.RoboPortfolioAsset, latestAssetPricess map[string]float64, totalValue float64, totalCash *float64) (float64, error) {
	totalSellAmount := 0.0
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(portfolioCategories))

	for _, asset := range portfolioCategories {
		wg.Add(1)
		go func(asset *models.RoboPortfolioAsset) {
			defer wg.Done()

			latestPrice, exists := latestAssetPricess[asset.Symbol]
			if !exists {
				errCh <- fmt.Errorf("latest price for asset %s not found", asset.Symbol)
				return
			}

			curValue := latestPrice * asset.SharesOwned
			targetValue := totalValue * asset.Percentage / 100
			amountToSell := curValue - targetValue

			sharesToSell := amountToSell / latestPrice
			asset.SharesOwned -= sharesToSell
			asset.TotalInvested -= amountToSell
			if asset.TotalInvested < 0 {
				asset.TotalInvested = 0
			}
			if asset.SharesOwned > 0 {
				asset.AvgBuyPrice = asset.TotalInvested / asset.SharesOwned
			}

			mu.Lock()
			totalSellAmount += amountToSell
			*totalCash += amountToSell
			mu.Unlock()

			transaction := &models.RoboPortfolioTransaction{
				RoboPortfolioID: asset.RoboPortfolioCategoryID,
				TransactionType: "rebalance:sell",
				TotalAmount:     amountToSell,
				Symbol:          &asset.Symbol,
				Name:            &asset.Name,
				Price:           &latestPrice,
				SharesAmount:    &sharesToSell,
			}
			if err := s.repo.CreateRoboPortfolioTransaction(transaction); err != nil {
				errCh <- fmt.Errorf("failed to create transaction for asset %s: %w", asset.Symbol, err)
			}
		}(asset)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		return 0, err
	}
	return totalSellAmount, nil
}

func (s *roboPortfolioServiceImpl) buyUnderPerformingAssets(portfolioCategories []*models.RoboPortfolioAsset, latestAssetPricess map[string]float64, totalValue float64, totalCash *float64, failReason *string) (float64, error) {
	totalBuyAmount := 0.0
	var mu sync.Mutex
	var once sync.Once
	var wg sync.WaitGroup
	errCh := make(chan error, len(portfolioCategories))

	for _, asset := range portfolioCategories {
		wg.Add(1)
		go func(asset *models.RoboPortfolioAsset) {
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
				once.Do(func() {
					*failReason = "There is not enough funds when rebalancing assets"

				})
				amountToBuy = *totalCash
			}
			*totalCash -= amountToBuy
			totalBuyAmount += amountToBuy
			mu.Unlock()

			if amountToBuy == 0 {
				return
			}
			sharesToBuy := amountToBuy / latestPrice
			asset.SharesOwned += sharesToBuy
			asset.TotalInvested += amountToBuy
			if asset.SharesOwned > 0 {
				asset.AvgBuyPrice = asset.TotalInvested / asset.SharesOwned
			}

			transaction := &models.RoboPortfolioTransaction{
				RoboPortfolioID: asset.RoboPortfolioCategoryID,
				TransactionType: "rebalance:buy",
				TotalAmount:     amountToBuy,
				Symbol:          &asset.Symbol,
				Name:            &asset.Name,
				Price:           &latestPrice,
				SharesAmount:    &sharesToBuy,
			}

			if err := s.repo.CreateRoboPortfolioTransaction(transaction); err != nil {
				errCh <- fmt.Errorf("failed to create transaction for asset %s: %w", asset.Symbol, err)
			}
		}(asset)
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		return 0, err
	}
	return totalBuyAmount, nil
}
