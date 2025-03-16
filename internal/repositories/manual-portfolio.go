package repositories

import (
	"errors"
	"fmt"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"gorm.io/gorm"
)

type ManualPortfolioRepo interface {
	GetManualPortfolios(userID uint) ([]*models.ManualPortfolio, error)
	GetManualPortfolio(userID uint, portfolioName string) (*models.ManualPortfolio, error)
	CreateManualPortfolio(portfolio *models.ManualPortfolio) error
	UpdateManualPortfolioName(portfolio *models.ManualPortfolio, newName string) error
	UpdateManualPortfolio(portfolio *models.ManualPortfolio) error
	DeleteManualPortfolio(portfolio *models.ManualPortfolio) error
}

type postgresManualPortfolioRepo struct {
	db *gorm.DB
}

func NewPostgresManualPortfolioRepo(db *gorm.DB) *postgresManualPortfolioRepo {
	return &postgresManualPortfolioRepo{db: db}
}

func (r *postgresManualPortfolioRepo) CreateManualPortfolio(portfolio *models.ManualPortfolio) error {
	if portfolio == nil {
		return commons.ErrNil
	}

	if err := r.db.Create(&portfolio).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}
		return err
	}
	fmt.Println("portfolio created", portfolio)
	return nil
}

func (r *postgresManualPortfolioRepo) GetManualPortfolios(userID uint) ([]*models.ManualPortfolio, error) {
	var portfolios []*models.ManualPortfolio
	if err := r.db.Where("user_id = ?", userID).Find(&portfolios).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return portfolios, nil
}

func (r *postgresManualPortfolioRepo) GetManualPortfolio(userID uint, portfolioName string) (*models.ManualPortfolio, error) {
	var portfolio models.ManualPortfolio
	if err := r.db.Where("user_id = ? AND name = ?", userID, portfolioName).Preload("Assets").First(&portfolio).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &portfolio, nil
}

func (r *postgresManualPortfolioRepo) UpdateManualPortfolioName(portfolio *models.ManualPortfolio, newName string) error {
	if portfolio == nil {
		return commons.ErrNil
	}
	portfolio.Name = newName
	if err := r.db.Save(&portfolio).Error; err != nil {
		return fmt.Errorf("failed to update portfolio name: %w", err)
	}
	return nil
}

func (r *postgresManualPortfolioRepo) UpdateManualPortfolio(portfolio *models.ManualPortfolio) error {
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

	for _, asset := range portfolio.Assets {
		if err := tx.Save(&asset).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update asset %s: %w", asset.Symbol, err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *postgresManualPortfolioRepo) DeleteManualPortfolio(portfolio *models.ManualPortfolio) error {

	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, asset := range portfolio.Assets {
		if err := tx.Delete(&asset).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete asset %s: %w", asset.Symbol, err)
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
