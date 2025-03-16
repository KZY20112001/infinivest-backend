package models

import "gorm.io/gorm"

type RoboPortfolio struct {
	gorm.Model
	UserID        uint                     `gorm:"primaryKey"`
	Categories    []*RoboPortfolioCategory `json:"categories"`
	RebalanceFreq *string                  `json:"rebalanceFreq"`
}

type RoboPortfolioCategory struct {
	gorm.Model
	RoboPortfolioID     uint                  `gorm:"not null"`
	RoboPortfolioUserID uint                  `gorm:"not null"`
	Name                string                `json:"name"`
	TotalPercentage     float64               `json:"totalPercentage"`
	TotalAmount         float64               `json:"totalAmount"`
	Assets              []*RoboPortfolioAsset `json:"assets"`
}

type RoboPortfolioAsset struct {
	gorm.Model
	RoboPortfolioCategoryID uint

	Symbol     string  `json:"symbol"`
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`

	SharesOwned   float64 `json:"sharesOwned"`
	TotalInvested float64 `json:"totalInvested"`
	AvgBuyPrice   float64 `json:"avgBuyPrice"`
}
