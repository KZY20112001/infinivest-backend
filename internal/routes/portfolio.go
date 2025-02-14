package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterPortfolioRoutes(r *gin.Engine, h *handlers.PortfolioHandler) {
	portfolioGroup := r.Group("/portfolio")
	portfolioGroup.Use(middlewares.AuthMiddleware())

	// for testing only
	portfolioGroup.GET("/:user_id/:portfolio_id", h.GetPortfolio)

	roboAdvisorGroup := portfolioGroup.Group("/robo-advisor")
	{
		roboAdvisorGroup.POST("/generate/categories", h.GenerateRoboAdvisorPortfolio)
		roboAdvisorGroup.POST("/generate/assets", h.GenerateAssetAllocation)
		roboAdvisorGroup.POST("/confirm", h.ConfirmGeneratedRoboPortfolio)
		roboAdvisorGroup.GET("/", h.GetRoboPortfolio)
		roboAdvisorGroup.POST("/", h.AddMoneyToRoboPortfolio)

	}

	manualGroup := portfolioGroup.Group("/manual")
	{
		manualGroup.POST("/")
		manualGroup.PATCH("/")
		manualGroup.GET("/", h.GetManualPortfolios)
	}
}
