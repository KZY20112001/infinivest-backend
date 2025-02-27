package models

import "gorm.io/gorm"

type RoboPortfolio struct {
	gorm.Model
	UserID        uint `gorm:"primaryKey"`
	Name          string
	Category      []*RoboPortfolioCategory
	RebalanceFreq *string
}

type RoboPortfolioCategory struct {
	gorm.Model
	RoboPortfolioID     uint `gorm:"not null"`
	RoboPortfolioUserID uint `gorm:"not null"`
	Name                string
	TotalPercentage     float64
	TotalAmount         float64
	Assets              []*RoboPortfolioAsset
}

type RoboPortfolioAsset struct {
	gorm.Model
	RoboPortfolioCategoryID uint

	Symbol     string
	Percentage float64

	SharesOwned   float64
	TotalInvested float64
	AvgBuyPrice   float64
}
