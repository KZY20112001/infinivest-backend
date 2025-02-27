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

func (h *ManualPortfolioHandler) GetManualPortfolios(c *gin.Context) {
	userID := c.GetUint("id")
	portfolios, err := h.service.GetManualPortfolios(userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolios": portfolios})
}

func (h *ManualPortfolioHandler) GetManualPortfolio(c *gin.Context) {
	portfolioName := c.Param("portfolio_name")
	userID := c.GetUint("id")
	portfolio, err := h.service.GetManualPortfolio(userID, portfolioName)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}
