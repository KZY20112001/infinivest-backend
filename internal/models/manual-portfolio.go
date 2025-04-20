package models

import "gorm.io/gorm"

type ManualPortfolio struct {
	gorm.Model
	UserID       uint                          `gorm:"not null;index:idx_user_name,unique"`
	Name         string                        `gorm:"index:idx_user_name,unique" json:"name"`
	Assets       []*ManualPortfolioAsset       `json:"assets"`
	TotalCash    float64                       `json:"totalCash"`
	Transactions []*ManualPortfolioTransaction `json:"manualPortfolioTransactions"`
}

type ManualPortfolioAsset struct {
	gorm.Model
	ManualPortfolioID     uint    `gorm:"not null;index"`
	ManualPortfolioUserID uint    `gorm:"not null;index"`
	Symbol                string  `json:"symbol"`
	Name                  string  `json:"name"`
	SharesOwned           float64 `json:"sharesOwned"`
	TotalInvested         float64 `json:"totalInvested"`
	AvgBuyPrice           float64 `json:"avgBuyPrice"`
}
