package repositories

import (
	"errors"
	"fmt"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"gorm.io/gorm"
)

type RoboPortfolioRepo interface {
	CreateRoboPortfolio(portfolio *models.RoboPortfolio) error
	GetRoboPortfolioDetails(userID uint) (*models.RoboPortfolio, error)
	UpdateRebalanceFreq(userID uint, freq string) error
	UpdateRoboPortfolio(portfolio *models.RoboPortfolio) error
	DeleteRoboPortfolio(portfolio *models.RoboPortfolio) error
}

type postgresRoboPortfolioRepo struct {
	db *gorm.DB
}

func NewPostgresRoboPortfolioRepo(db *gorm.DB) *postgresRoboPortfolioRepo {
	return &postgresRoboPortfolioRepo{db: db}
}

func (r *postgresRoboPortfolioRepo) CreateRoboPortfolio(portfolio *models.RoboPortfolio) error {
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

func (r *postgresRoboPortfolioRepo) GetRoboPortfolioDetails(userID uint) (*models.RoboPortfolio, error) {
	var portfolio models.RoboPortfolio
	if err := r.db.
		Where("user_id = ?", userID).
		Preload("Categories").
		Preload("Categories.Assets").
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

func (r *postgresRoboPortfolioRepo) UpdateRoboPortfolio(portfolio *models.RoboPortfolio) error {
	if portfolio == nil {
		return commons.ErrNil
	}

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

	for _, category := range portfolio.Categories {
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

func (r *postgresRoboPortfolioRepo) UpdateRebalanceFreq(userID uint, freq string) error {
	portfolio, err := r.GetRoboPortfolioDetails(userID)
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

func (r *postgresRoboPortfolioRepo) DeleteRoboPortfolio(portfolio *models.RoboPortfolio) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, category := range portfolio.Categories {
		for _, asset := range category.Assets {
			if err := tx.Delete(&asset).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to delete asset %s: %w", asset.Symbol, err)
			}
		}
	}

	for _, category := range portfolio.Categories {
		if err := tx.Delete(&category).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete category %s: %w", category.Name, err)
		}
	}

	if err := tx.Delete(&portfolio).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete portfolio: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
