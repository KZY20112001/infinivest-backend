package models

import "gorm.io/gorm"

type RoboPortfolioTransaction struct {
	gorm.Model
	RoboPortfolioID uint

	TransactionType string  `json:"transactionType"` // "buy" or "sell" or "dividend" or "deposit" or "withdrawal" or or "rebalance:sell" or "rebalance:buy"
	TotalAmount     float64 `json:"totalAmount"`

	Symbol       *string  `json:"symbol"`
	Name         *string  `json:"name"`
	Price        *float64 `json:"price"`
	SharesAmount *float64 `json:"sharesAmount"`
}

type ManualPortfolioTransaction struct {
	gorm.Model
	ManualPortfolioUserID uint `gorm:"not null;index"`
	ManualPortfolioID     uint `gorm:"not null;index"`

	TransactionType string  `json:"transactionType"` // "buy" or "sell" or "dividend" or "deposit" or "withdrawal"
	TotalAmount     float64 `json:"totalAmount"`

	Symbol       *string  `json:"symbol"`
	Name         *string  `json:"name"`
	Price        *float64 `json:"price"`
	SharesAmount *float64 `json:"sharesAmount"`
}
