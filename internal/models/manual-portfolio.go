package models

import "gorm.io/gorm"

type ManualPortfolio struct {
	gorm.Model
	UserID    uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	Assets    []*ManualPortfolioAssets
	TotalCash float64
}

type ManualPortfolioAssets struct {
	gorm.Model
	ManualPortfolioID     uint
	ManualPortfolioUserID uint
	Symbol                string

	SharesOwned   float64
	TotalInvested float64
	AvgBuyPrice   float64
}
