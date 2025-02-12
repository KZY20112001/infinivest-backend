package repositories

import (
	"errors"

	"github.com/KZY20112001/infinivest-backend/internal/constants"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"gorm.io/gorm"
)

type PortfolioRepo interface {
	CreatePortfolio(portfolio *models.Portfolio) error
	GetPortfolio(portfolioID uint, userID uint) (*models.Portfolio, error)
	GetRoboPortfolio(userID uint) (*models.Portfolio, error)
	GetManualPortfolios(userID uint) ([]models.Portfolio, error)
}

type postgresPortfolioRepo struct {
	db *gorm.DB
}

func NewPostgresPortfolioRepo(db *gorm.DB) *postgresPortfolioRepo {
	return &postgresPortfolioRepo{db: db}
}

func (ptr *postgresPortfolioRepo) CreatePortfolio(portfolio *models.Portfolio) error {
	if portfolio == nil {
		return constants.ErrNil
	}
	if err := ptr.db.Create(&portfolio).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return constants.ErrDuplicate
		}
		return err
	}

	return nil
}

func (ptr *postgresPortfolioRepo) GetPortfolio(portfolioID uint, userID uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio

	if err := ptr.db.
		Where("id = ? AND user_id = ?", portfolioID, userID).
		Preload("Category").
		Preload("Category.Assets").
		First(&portfolio).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return &portfolio, nil
}

func (ptr *postgresPortfolioRepo) GetRoboPortfolio(userID uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	if err := ptr.db.
		Where("user_id = ? AND is_robo_advisor = ?", userID, true).
		First(&portfolio).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return &portfolio, nil
}

func (ptr *postgresPortfolioRepo) GetManualPortfolios(userID uint) ([]models.Portfolio, error) {
	var portfolios []models.Portfolio
	if err := ptr.db.
		Where("user_id = ? AND is_robo_advisor = ?", userID, false).
		Find(&portfolios).Error; err != nil {
		return nil, err
	}
	return portfolios, nil
}
