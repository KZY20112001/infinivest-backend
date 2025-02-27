package models

import "gorm.io/gorm"

type ManualPortfolio struct {
	gorm.Model
	UserID    uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	Assets    []*ManualPortfolioAsset
	TotalCash float64
}

type ManualPortfolioAsset struct {
	gorm.Model
	ManualPortfolioID     uint
	ManualPortfolioUserID uint
	Symbol                string

	SharesOwned   float64
	TotalInvested float64
	AvgBuyPrice   float64
}
