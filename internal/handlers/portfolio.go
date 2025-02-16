package handlers

import (
	"fmt"
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	portfolioService services.PortfolioService
	genAIService     services.GenAIService
}

func NewPortfolioHandler(ps services.PortfolioService, gs services.GenAIService) *PortfolioHandler {
	return &PortfolioHandler{portfolioService: ps, genAIService: gs}
}

func (h *PortfolioHandler) GenerateRoboAdvisorPortfolio(c *gin.Context) {
	bankName := c.PostForm("bank_name")
	riskToleranceLevel := c.PostForm("risk_tolerance_level")
	if riskToleranceLevel == "" {
		HandleError(c, fmt.Errorf("risk tolerance level is required"))
		return
	}

	bankStatement, err := c.FormFile("bank_statement")
	if err != nil {
		HandleError(c, err)
		return
	}

	recommendation, err := h.genAIService.GenerateRoboAdvisorPortfolio(bankStatement, bankName, riskToleranceLevel)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, recommendation)
}

func (h *PortfolioHandler) GenerateAssetAllocation(c *gin.Context) {
	var req dto.AssetAllocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetAllocations, err := h.genAIService.GenerateAssetAllocations(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assetAllocations)
}

func (h *PortfolioHandler) ConfirmGeneratedRoboPortfolio(c *gin.Context) {
	var req dto.ConfirmPortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("id")
	err := h.portfolioService.ConfirmGeneratedRoboPortfolio(req, userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created the portfolio"})
}

func (h *PortfolioHandler) GetRoboPortfolio(c *gin.Context) {
	userID := c.GetUint("id")
	portfolio, err := h.portfolioService.GetRoboPortfolio(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}

func (h *PortfolioHandler) AddMoneyToRoboPortfolio(c *gin.Context) {
	var req dto.AddMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("id")
	portfolio, err := h.portfolioService.AddMoneyToRoboPortfolio(userID, req.Amount)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}

func (h *PortfolioHandler) UpdateRebalanceFreq(c *gin.Context) {
	var req dto.UpdateRebalanceFreqRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("id")
	err := h.portfolioService.UpdateRebalanceFreq(userID, req.Frequency)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated the rebalance frequency"})
}

// handlers for manual-portfolios
func (h *PortfolioHandler) GetManualPortfolios(c *gin.Context) {
	userID := c.GetUint("id")
	portfolios, err := h.portfolioService.GetManualPortfolios(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolios": portfolios})
}
