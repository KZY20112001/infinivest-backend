package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterPortfolioRoutes(r *gin.Engine, rh *handlers.RoboPortfolioHandler, mh *handlers.ManualPortfolioHandler) {
	portfolioGroup := r.Group("/portfolio")
	portfolioGroup.Use(middlewares.AuthMiddleware())

	roboAdvisorGroup := portfolioGroup.Group("/robo-advisor")
	{
		roboAdvisorGroup.POST("/generate/categories", rh.GenerateRoboAdvisorPortfolio)
		roboAdvisorGroup.POST("/generate/assets", rh.GenerateAssetAllocation)
		roboAdvisorGroup.POST("/confirm", rh.ConfirmGeneratedRoboPortfolio)
		roboAdvisorGroup.GET("/", rh.GetRoboPortfolio)
		roboAdvisorGroup.POST("/", rh.AddMoneyToRoboPortfolio)
		roboAdvisorGroup.PUT("/rebalance-freq", rh.UpdateRebalanceFreq)
	}

	manualGroup := portfolioGroup.Group("/manual")
	{
		manualGroup.POST("/")
		manualGroup.PATCH("/")
		manualGroup.GET("/", mh.GetManualPortfolios)
	}
}
