package models

import "gorm.io/gorm"

type Portfolio struct {
	gorm.Model
	UserID        uint `gorm:"primaryKey"`
	Name          string
	IsRoboAdvisor bool
	Category      []*PortfolioCategory
	RebalanceFreq *string
}

type PortfolioCategory struct {
	gorm.Model
	PortfolioID     uint `gorm:"not null"`
	PortfolioUserID uint `gorm:"not null"`
	Name            string
	TotalPercentage float64
	TotalAmount     float64
	Assets          []*PortfolioAsset
}

type PortfolioAsset struct {
	gorm.Model
	PortfolioCategoryID uint

	Symbol     string
	Percentage float64

	SharesOwned   float64
	TotalInvested float64
	AvgBuyPrice   float64
}
