package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterPortfolioRoutes(r *gin.Engine, h *handlers.PortfolioHandler) {
	portfolioGroup := r.Group("/portfolio")
	portfolioGroup.Use(middlewares.AuthMiddleware())

	roboAdvisorGroup := portfolioGroup.Group("/robo-advisor")
	{
		roboAdvisorGroup.POST("/generate/categories", h.GenerateRoboAdvisorPortfolio)
		roboAdvisorGroup.POST("/generate/assets", h.GenerateAssetAllocation)
		roboAdvisorGroup.POST("/confirm", h.ConfirmGeneratedRoboPortfolio)
		roboAdvisorGroup.GET("/", h.GetRoboPortfolio)
		roboAdvisorGroup.POST("/", h.AddMoneyToRoboPortfolio)
		roboAdvisorGroup.PUT("/rebalance-freq", h.UpdateRebalanceFreq)
	}

	manualGroup := portfolioGroup.Group("/manual")
	{
		manualGroup.POST("/")
		manualGroup.PATCH("/")
		manualGroup.GET("/", h.GetManualPortfolios)
	}
}
