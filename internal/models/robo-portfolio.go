package models

import "gorm.io/gorm"

type RoboPortfolio struct {
	gorm.Model
	UserID          uint                        `gorm:"uniqueIndex;not null"`
	Categories      []*RoboPortfolioCategory    `json:"categories"`
	RebalanceFreq   *string                     `json:"rebalanceFreq"`
	RebalanceEvents []*AutoRebalanceEvent       `json:"rebalanceEvents"`
	Transactions    []*RoboPortfolioTransaction `json:"roboPortfolioTransactions"`
}

type RoboPortfolioCategory struct {
	gorm.Model
	RoboPortfolioID uint                  `gorm:"not null"`
	Name            string                `json:"name"`
	TotalPercentage float64               `json:"totalPercentage"`
	TotalAmount     float64               `json:"totalAmount"`
	Assets          []*RoboPortfolioAsset `json:"assets"`
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

type AutoRebalanceEvent struct {
	gorm.Model
	RoboPortfolioID uint `gorm:"not null"`

	Transactions    []*RoboPortfolioTransaction `json:"roboPortfolioTransactions"`
	TotalBuyAmount  float64                     `json:"totalBuyAmount"`
	TotalSellAmount float64                     `json:"totalSellAmount"`
	NetChange       float64                     `json:"netChange"`

	Success bool    `json:"success"` // was this rebalance completed?
	Reason  *string `json:"reason"`  // optional error or reason (e.g. "insufficient funds")

	PortfolioValueBefore float64 `json:"portfolioValueBefore"`
	PortfolioValueAfter  float64 `json:"portfolioValueAfter"`
	GainOrLoss           float64 `json:"gainOrLoss"`
}
