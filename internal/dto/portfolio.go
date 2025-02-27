package dto

type RoboAdvisorPortfolio struct {
	LargeCapBlend       float64 `json:"large_cap_blend"`
	SmallCapBlend       float64 `json:"small_cap_blend"`
	InternationalStocks float64 `json:"international_stocks"`
	EmergingMarkets     float64 `json:"emerging_markets"`
	IntermediateBonds   float64 `json:"intermediate_bonds"`
	InternationalBonds  float64 `json:"international_bonds"`
	Cash                float64 `json:"cash"`
}

type RoboAdvisorRecommendationResponse struct {
	Portfolio RoboAdvisorPortfolio `json:"portfolio"`
	Reason    string               `json:"reason"`
}

type AssetAllocationRequest struct {
	Portfolio map[string]float64 `json:"portfolio"`
}

type Asset struct {
	Symbol     string  `json:"symbol"`
	Percentage float64 `json:"percentage"`
}

type Assets struct {
	Assets []Asset `json:"assets"`
}

type AssetAllocationResponse struct {
	Allocations map[string]Assets `json:"allocations"`
}

type ConfirmPortfolioRequest struct {
	Portfolio   map[string]float64 `json:"portfolio"`
	Allocations map[string]Assets  `json:"allocations"`
	Frequency   string             `json:"frequency"`
}

type UpdateRebalanceFreqRequest struct {
	Frequency string `json:"frequency"`
}

type AddMoneyRequest struct {
	Amount float64 `json:"amount"`
}

type WithdrawMoneyRequest struct {
	Amount float64 `json:"amount"`
}

type ManualPortfolioRequest struct {
	Name string `json:"name"`
}
