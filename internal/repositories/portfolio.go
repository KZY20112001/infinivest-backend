package repositories

import (
	"errors"
	"fmt"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"gorm.io/gorm"
)

type PortfolioRepo interface {
	GetPortfolio(portfolioID, userID uint) (*models.Portfolio, error)
	CreatePortfolio(portfolio *models.Portfolio) error
	GetRoboPortfolio(userID uint) (*models.Portfolio, error)
	UpdateRebalanceFreq(userID uint, freq string) error
	GetManualPortfolios(userID uint) ([]models.Portfolio, error)
	UpdatePortfolio(portfolio *models.Portfolio) error
}

type postgresPortfolioRepo struct {
	db *gorm.DB
}

func NewPostgresPortfolioRepo(db *gorm.DB) *postgresPortfolioRepo {
	return &postgresPortfolioRepo{db: db}
}

func (r *postgresPortfolioRepo) GetPortfolio(portfolioID, userID uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	if err := r.db.
		Where("id = ? AND user_id = ?", portfolioID, userID).
		Preload("Category").
		Preload("Category.Assets").
		First(&portfolio).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &portfolio, nil
}

func (r *postgresPortfolioRepo) CreatePortfolio(portfolio *models.Portfolio) error {
	if portfolio == nil {
		return commons.ErrNil
	}
	if err := r.db.Create(&portfolio).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
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
			return nil, gorm.ErrRecordNotFound
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, gorm.ErrDuplicatedKey
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

func (r *postgresPortfolioRepo) UpdatePortfolio(portfolio *models.Portfolio) error {
	tx := r.db.Begin()

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

func (r *postgresPortfolioRepo) UpdateRebalanceFreq(userID uint, freq string) error {
	portfolio, err := r.GetRoboPortfolio(userID)
	if err != nil {
		return err
	}
	if err := r.db.Model(&portfolio).Update("rebalance_freq", freq).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}
		return err
	}
	return nil
}
