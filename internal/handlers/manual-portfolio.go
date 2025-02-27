package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ManualPortfolioHandler struct {
	service services.ManualPortfolioService
}

func NewManualPortfolioHandler(ps services.ManualPortfolioService) *ManualPortfolioHandler {
	return &ManualPortfolioHandler{service: ps}
}

// handlers for manual-portfolios
func (h *ManualPortfolioHandler) GetManualPortfolios(c *gin.Context) {
	userID := c.GetUint("id")
	portfolios, err := h.service.GetManualPortfolios(userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolios": portfolios})
}
