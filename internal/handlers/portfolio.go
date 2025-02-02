package handlers

import (
	"fmt"
	"net/http"

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
	bankName := c.Param("bank_name")
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
	fmt.Println("name: ", bankStatement.Filename)
	recommendation, err := h.portfolioService.GenerateRoboAdvisorPortfolio(bankStatement, bankName, riskToleranceLevel)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, recommendation)
}
