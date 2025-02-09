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
}

func NewPortfolioHandler(ps services.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{portfolioService: ps}
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

	recommendation, err := h.portfolioService.GenerateRoboAdvisorPortfolio(bankStatement, bankName, riskToleranceLevel)
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
	assetAllocations, err := h.portfolioService.GenerateAssetAllocations(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assetAllocations)
}
