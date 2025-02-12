package models

import "gorm.io/gorm"

type Portfolio struct {
	gorm.Model
	UserID        uint `gorm:"primaryKey"`
	Name          string
	IsRoboAdvisor bool
	Category      []PortfolioCategory
}

type PortfolioCategory struct {
	ID              uint `gorm:"primaryKey"`
	PortfolioID     uint `gorm:"not null"`
	PortfolioUserID uint `gorm:"not null"`
	Name            string
	TotalPercentage float64
	Assets          []PortfolioAsset
}

type PortfolioAsset struct {
	ID                  uint `gorm:"primaryKey"`
	PortfolioCategoryID uint `gorm:"not null"`
	Symbol              string
	Percentage          float64
}
