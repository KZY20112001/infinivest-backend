package services

import (
	"strconv"

	"github.com/KZY20112001/infinivest-backend/internal/constants"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type PortfolioService interface {
	ConfirmGeneratedRoboPortfolio(req dto.ConfirmPortfolioRequest, userID uint) error
	GetRoboPortfolio(userID uint) (*models.Portfolio, error)
	GetManualPortfolios(userID uint) ([]models.Portfolio, error)
	GetPortfolio(portfolioID, userID uint) (*models.Portfolio, error)
}

type portfolioServiceImpl struct {
	repo           repositories.PortfolioRepo
	profileService ProfileService
}

func NewPortfolioService(pr repositories.PortfolioRepo, ps ProfileService) *portfolioServiceImpl {
	return &portfolioServiceImpl{repo: pr, profileService: ps}
}

func (s *portfolioServiceImpl) ConfirmGeneratedRoboPortfolio(req dto.ConfirmPortfolioRequest, userID uint) error {
	_, err := s.repo.GetRoboPortfolio(userID)
	if err == nil {
		return constants.ErrDuplicate
	}
	portfolio := models.Portfolio{
		UserID:        userID,
		Name:          strconv.FormatUint(uint64(userID), 10) + " Robo Advisor Portfolio",
		IsRoboAdvisor: true,
		Category:      []models.PortfolioCategory{},
	}

	for categoryName, assets := range req.Allocations {
		category := models.PortfolioCategory{
			PortfolioUserID: userID,
			PortfolioID:     portfolio.ID,
			Name:            categoryName,
			TotalPercentage: req.Portfolio[categoryName],
			Assets:          []models.PortfolioAsset{},
		}

		for _, asset := range assets.Assets {
			category.Assets = append(category.Assets, models.PortfolioAsset{
				PortfolioCategoryID: category.ID,
				Symbol:              asset.Symbol,
				Percentage:          asset.Percentage,
			})
		}

		portfolio.Category = append(portfolio.Category, category)
	}
	return s.repo.CreatePortfolio(&portfolio)
}

func (s *portfolioServiceImpl) GetRoboPortfolio(userID uint) (*models.Portfolio, error) {
	return s.repo.GetRoboPortfolio(userID)
}

func (s *portfolioServiceImpl) GetManualPortfolios(userID uint) ([]models.Portfolio, error) {
	return s.repo.GetManualPortfolios(userID)
}

func (s *portfolioServiceImpl) GetPortfolio(portfolioID, userID uint) (*models.Portfolio, error) {
	return s.repo.GetPortfolio(portfolioID, userID)
}
