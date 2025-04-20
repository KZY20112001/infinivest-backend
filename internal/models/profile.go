package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID            uint
	User              User
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	ProfileUrl        string  `json:"profileUrl"`
	ProfileID         string  `json:"profileID"`
	RiskTolerance     string  `json:"riskTolerance"`
	InvestmentStyle   string  `json:"investmentStyle"`
	InvestmentHorizon string  `json:"investmentHorizon"`
	AnnualIncome      float64 `json:"annualIncome"`
	ExperienceLevel   string  `json:"experienceLevel"`
}
