package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ManualPortfolioHandler struct {
	service services.ManualPortfolioService
}

func NewManualPortfolioHandler(ps services.ManualPortfolioService) *ManualPortfolioHandler {
	return &ManualPortfolioHandler{service: ps}
}

func (h *ManualPortfolioHandler) GetManualPortfoliosDetails(c *gin.Context) {
	userID := c.GetUint("id")
	portfolios, err := h.service.GetManualPortfoliosDetails(userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolios": portfolios})
}

func (h *ManualPortfolioHandler) GetManualPortfoliosSummaries(c *gin.Context) {
	userID := c.GetUint("id")
	portfolios, err := h.service.GetManualPortfoliosSummaries(userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolios": portfolios})

}

func (h *ManualPortfolioHandler) GetManualPortfolio(c *gin.Context) {
	portfolioName := c.Param("name")
	userID := c.GetUint("id")
	portfolio, err := h.service.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}

func (h *ManualPortfolioHandler) GetPortfolioValue(c *gin.Context) {
	portfolioName := c.Param("name")
	userID := c.GetUint("id")
	portfolioValue, err := h.service.GetPorfolioValue(userID, portfolioName)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"amount": portfolioValue})
}

func (h *ManualPortfolioHandler) CreateManualPortfolio(c *gin.Context) {
	var req dto.ManualPortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")
	if err := h.service.CreateManualPortfolio(userID, req.Name); err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Portfolio created successfully"})
}

func (h *ManualPortfolioHandler) UpdatePortfolioName(c *gin.Context) {
	var req dto.ManualPortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")
	portfolioName := c.Param("name")
	if err := h.service.UpdatePortfolioName(userID, portfolioName, req.Name); err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Portfolio name updated successfully"})
}

func (h *ManualPortfolioHandler) AddMoneyToManualPortfolio(c *gin.Context) {
	var req dto.AddMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")
	portfolioName := c.Param("name")
	if err := h.service.AddMoneyToManualPortfolio(userID, portfolioName, req.Amount); err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Money added successfully"})
}

func (h *ManualPortfolioHandler) WithDrawMoneyFromManualPortfolio(c *gin.Context) {
	var req dto.WithdrawMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")
	portfolioName := c.Param("name")
	amountWithdrawn, err := h.service.WithdrawMoneyFromManualPortfolio(userID, portfolioName, req.Amount)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"amount": amountWithdrawn})
}

func (h *ManualPortfolioHandler) BuyAssetForManualPortfolio(c *gin.Context) {
	var req dto.ManualPortfolioBuyAssetRequest
	userID := c.GetUint("id")
	portfolioName := c.Param("name")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.BuyAssetForManualPortfolio(userID, portfolioName, req.Name, req.Symbol, req.SharesAmount); err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shares bought successfully"})
}

func (h *ManualPortfolioHandler) SellAssetForManualPortfolio(c *gin.Context) {
	var req dto.ManualPortfolioSellAssetRequest
	userID := c.GetUint("id")
	portfolioName := c.Param("name")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SellAssetForManualPortfolio(userID, portfolioName, req.Symbol, req.SharesAmount); err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shares sold successfully"})
}

func (h *ManualPortfolioHandler) DeleteManualPortfolio(c *gin.Context) {
	portfolioName := c.Param("name")
	userID := c.GetUint("id")
	if err := h.service.DeleteManualPortfolio(userID, portfolioName); err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Portfolio deleted successfully"})

}
