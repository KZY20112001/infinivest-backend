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
		roboAdvisorGroup.POST("/add", rh.AddMoneyToRoboPortfolio)
		roboAdvisorGroup.POST("/withdraw", rh.WithDrawMoneyFromRoboPortfolio)
		roboAdvisorGroup.PUT("/rebalance-freq", rh.UpdateRebalanceFreq)
		roboAdvisorGroup.DELETE("/", rh.DeleteRoboPortfolio)
	}

	manualGroup := portfolioGroup.Group("/manual")
	{
		manualGroup.POST("/", mh.CreateManualPortfolio)
		manualGroup.POST("/:name/add", mh.AddMoneyToManualPortfolio)
		manualGroup.POST("/:name/withdraw", mh.WithDrawMoneyFromManualPortfolio)
		manualGroup.GET("/", mh.GetManualPortfolios)
		manualGroup.GET("/:name", mh.GetManualPortfolio)
	}
}
