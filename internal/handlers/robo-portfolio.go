package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type RoboPortfolioHandler struct {
	service      services.RoboPortfolioService
	genAIService services.GenAIService
}

func NewRoboPortfolioHandler(ps services.RoboPortfolioService, gs services.GenAIService) *RoboPortfolioHandler {
	return &RoboPortfolioHandler{service: ps, genAIService: gs}
}

func (h *RoboPortfolioHandler) GenerateRoboAdvisorPortfolio(c *gin.Context) {
	bankName := c.PostForm("bank_name")
	if bankName == "" {
		commons.HandleError(c, fmt.Errorf("bank name is required"))
		return
	}

	riskToleranceLevel := c.PostForm("risk_tolerance_level")
	if riskToleranceLevel == "" {
		commons.HandleError(c, fmt.Errorf("risk tolerance level is required"))
		return
	}

	bankStatement, err := c.FormFile("bank_statement")

	if err != nil {
		commons.HandleError(c, err)
		return
	}

	recommendation, err := h.genAIService.GenerateRoboAdvisorPortfolio(bankStatement, bankName, riskToleranceLevel)
	if err != nil {
		commons.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, recommendation)
}

func (h *RoboPortfolioHandler) GenerateAssetAllocation(c *gin.Context) {
	var req dto.AssetAllocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetAllocations, err := h.genAIService.GenerateAssetAllocations(req)
	if err != nil {
		commons.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assetAllocations)
}

func (h *RoboPortfolioHandler) ConfirmGeneratedRoboPortfolio(c *gin.Context) {
	var req dto.ConfirmPortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := commons.RebalancingThresholds[req.Frequency]; !exists {
		commons.HandleError(c, fmt.Errorf("invalid rebalance frequency"))
		return
	}
	userID := c.GetUint("id")
	err := h.service.ConfirmGeneratedRoboPortfolio(req, userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created the portfolio"})
}

func (h *RoboPortfolioHandler) GetRoboPortfolioDetails(c *gin.Context) {
	userID := c.GetUint("id")
	portfolio, err := h.service.GetRoboPortfolioDetails(userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}

func (h *RoboPortfolioHandler) GetRoboPortfolioSummary(c *gin.Context) {
	userID := c.GetUint("id")
	summary, err := h.service.GetRoboPortfolioSummary(userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

func (h *RoboPortfolioHandler) AddMoneyToRoboPortfolio(c *gin.Context) {
	var req dto.AddMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")
	portfolio, err := h.service.AddMoneyToRoboPortfolio(c.Request.Context(), userID, req.Amount)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}

func (h *RoboPortfolioHandler) WithDrawMoneyFromRoboPortfolio(c *gin.Context) {
	var req dto.WithdrawMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")
	amountWithdrawn, err := h.service.WithDrawMoneyFromRoboPortfolio(c.Request.Context(), userID, req.Amount)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"amount": amountWithdrawn})
}

func (h *RoboPortfolioHandler) UpdateRebalanceFreq(c *gin.Context) {
	var req dto.UpdateRebalanceFreqRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := commons.RebalancingThresholds[req.Frequency]; !exists {
		commons.HandleError(c, fmt.Errorf("invalid rebalance frequency"))
		return
	}
	userID := c.GetUint("id")
	err := h.service.UpdateRebalanceFreq(c.Request.Context(), userID, req.Frequency)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated the rebalance frequency"})
}

func (h *RoboPortfolioHandler) DeleteRoboPortfolio(c *gin.Context) {
	userID := c.GetUint("id")
	err := h.service.DeleteRoboPortfolio(c.Request.Context(), userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted the portfolio"})
}

func (h *RoboPortfolioHandler) GetRoboPortfolioTransactions(c *gin.Context) {
	userID := c.GetUint("id")
	limitStr := c.Query("limit")

	limit := 0
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		limit = parsedLimit
	}
	transactions, err := h.service.GetRoboPortfolioTransactions(userID, limit)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func (h *RoboPortfolioHandler) GetRebalanceEvents(c *gin.Context) {
	userID := c.GetUint("id")
	rebalanceDetails, err := h.service.GetRebalanceEvents(c.Request.Context(), userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"rebalance_details": rebalanceDetails})
}

func (h *RoboPortfolioHandler) UpdateLastSeenRebalanceEvent(c *gin.Context) {
	userID := c.GetUint("id")
	err := h.service.UpdateLastSeenRebalanceTime(c.Request.Context(), userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated the last seen rebalance event"})
}

func (h *RoboPortfolioHandler) RebalanceRoboPortfolio(c *gin.Context) {
	userID := c.GetUint("id")
	portfolio, err := h.service.GetRoboPortfolioDetails(userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	portfolio, err = h.service.RebalancePortfolio(c.Request.Context(), userID, portfolio.ID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": portfolio})

}
