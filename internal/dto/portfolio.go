package dto

type RoboAdvisorPortfolio struct {
	LargeCapBlend       float64 `json:"largeCapBlend"`
	SmallCapBlend       float64 `json:"smallCapBlend"`
	InternationalStocks float64 `json:"internationalStocks"`
	EmergingMarkets     float64 `json:"emergingMarkets"`
	IntermediateBonds   float64 `json:"intermediateBonds"`
	InternationalBonds  float64 `json:"internationalBonds"`
	Cash                float64 `json:"cash"`
}

type RoboAdvisorRecommendationResponse struct {
	Portfolio RoboAdvisorPortfolio `json:"portfolio"`
	Reason    string               `json:"reason"`
}

type RoboPortfolioSummaryResponse struct {
	RebalanceFreq string  `json:"rebalanceFreq"`
	TotalValue    float64 `json:"totalValue"`
	TotalInvested float64 `json:"totalInvested"`
}

type AssetAllocationRequest struct {
	Portfolio map[string]float64 `json:"portfolio"`
}

type Asset struct {
	Name       string  `json:"name"`
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

type ManualPortfolioBuyAssetRequest struct {
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	SharesAmount float64 `json:"sharesAmount"`
}

type ManualPortfolioSellAssetRequest struct {
	Symbol       string  `json:"symbol"`
	SharesAmount float64 `json:"sharesAmount"`
}

type ManualPortfolioSummaryResponse struct {
	TotalValue    float64 `json:"totalValue"`
	TotalInvested float64 `json:"totalInvested"`
	Name          string  `json:"name"`
}
