package repositories

import (
	"errors"
	"fmt"

	"github.com/KZY20112001/infinivest-backend/internal/constants"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"gorm.io/gorm"
)

type PortfolioRepo interface {
	CreatePortfolio(portfolio *models.Portfolio) error
	GetRoboPortfolio(userID uint) (*models.Portfolio, error)
	GetManualPortfolios(userID uint) ([]models.Portfolio, error)
	UpdatePortfolio(portfolio *models.Portfolio) error
	GetPortfolio(portfolioID uint, userID uint) (*models.Portfolio, error)
}

type postgresPortfolioRepo struct {
	db *gorm.DB
}

func NewPostgresPortfolioRepo(db *gorm.DB) *postgresPortfolioRepo {
	return &postgresPortfolioRepo{db: db}
}

func (r *postgresPortfolioRepo) CreatePortfolio(portfolio *models.Portfolio) error {
	if portfolio == nil {
		return constants.ErrNil
	}
	if err := r.db.Create(&portfolio).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return constants.ErrDuplicate
		}
		return err
	}

	return nil
}

func (r *postgresPortfolioRepo) GetRoboPortfolio(userID uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	if err := r.db.
		Where("user_id = ? AND is_robo_advisor = ?", userID, true).
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

func (r *postgresPortfolioRepo) GetManualPortfolios(userID uint) ([]models.Portfolio, error) {
	var portfolios []models.Portfolio
	if err := r.db.
		Where("user_id = ? AND is_robo_advisor = ?", userID, false).
		Find(&portfolios).Error; err != nil {
		return nil, err
	}
	return portfolios, nil
}

func (ptr *postgresPortfolioRepo) UpdatePortfolio(portfolio *models.Portfolio) error {
	tx := ptr.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(&portfolio).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update portfolio: %w", err)
	}

	for _, category := range portfolio.Category {
		if err := tx.Save(&category).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update category %s: %w", category.Name, err)
		}

		for _, asset := range category.Assets {
			if err := tx.Save(&asset).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update asset %s: %w", asset.Symbol, err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *postgresPortfolioRepo) GetPortfolio(portfolioID uint, userID uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio

	if err := r.db.
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
