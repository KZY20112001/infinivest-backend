package models

import "gorm.io/gorm"

type AssetDetails struct {
	Symbol       *string  `json:"symbol"`
	Name         *string  `json:"name"`
	Price        *float64 `json:"price"`
	SharesAmount *float64 `json:"sharesAmount"`
}

type TransactionDetails struct {
	TransactionType string  `json:"transactionType"` // "buy" or "sell" or "dividend" or "deposit" or "withdrawal"
	TotalAmount     float64 `json:"totalAmount"`
}

type RoboPortfolioTransaction struct {
	gorm.Model
	RoboPortfolioID      uint
	AutoRebalanceEventID *uint

	TransactionDetails
	AssetDetails
}

type ManualPortfolioTransaction struct {
	gorm.Model
	ManualPortfolioID uint

	TransactionDetails
	AssetDetails
}
