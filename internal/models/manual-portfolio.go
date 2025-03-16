package models

import "gorm.io/gorm"

type ManualPortfolio struct {
	gorm.Model
	UserID    uint                    `gorm:"primaryKey"`
	Name      string                  `gorm:"uniqueIndex" json:"name"`
	Assets    []*ManualPortfolioAsset `json:"assets"`
	TotalCash float64                 `json:"totalCash"`
}

type ManualPortfolioAsset struct {
	gorm.Model
	ManualPortfolioID     uint
	ManualPortfolioUserID uint
	Symbol                string  `json:"symbol"`
	Name                  string  `json:"name"`
	SharesOwned           float64 `json:"sharesOwned"`
	TotalInvested         float64 `json:"totalInvested"`
	AvgBuyPrice           float64 `json:"avgBuyPrice"`
}
