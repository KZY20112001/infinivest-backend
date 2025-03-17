package dto

type ProfileRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name,omitempty"`
	ProfileUrl string `json:"profile_url,omitempty"`
	ProfileID  string `json:"profile_id,omitempty"`

	RiskTolerance     string  `json:"risk_tolerance,omitempty"`
	InvestmentStyle   string  `json:"investment_style,omitempty"`
	InvestmentHorizon string  `json:"investment_horizon,omitempty"`
	AnnualIncome      float64 `json:"annual_income,omitempty"`
	ExperienceLevel   string  `json:"experience_level,omitempty"`
}
