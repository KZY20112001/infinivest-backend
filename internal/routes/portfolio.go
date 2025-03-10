package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterPortfolioRoutes(r *gin.Engine, rh *handlers.RoboPortfolioHandler, mh *handlers.ManualPortfolioHandler) {
	portfolioGroup := r.Group("/portfolio")
	portfolioGroup.Use(middlewares.AuthMiddleware())

	roboAdvisorGroup := portfolioGroup.Group("/robo-portfolio")
	{
		roboAdvisorGroup.GET("/details", rh.GetRoboPortfolio)
		roboAdvisorGroup.GET("/summary", rh.GetRoboPortfolioSummary)

		// TODO: Get transactions
		roboAdvisorGroup.GET("/transactions")

		roboAdvisorGroup.POST("/generate/categories", rh.GenerateRoboAdvisorPortfolio)
		roboAdvisorGroup.POST("/generate/assets", rh.GenerateAssetAllocation)
		roboAdvisorGroup.POST("/confirm", rh.ConfirmGeneratedRoboPortfolio)
		roboAdvisorGroup.POST("/add", rh.AddMoneyToRoboPortfolio)
		roboAdvisorGroup.POST("/withdraw", rh.WithDrawMoneyFromRoboPortfolio)
		roboAdvisorGroup.PUT("/rebalance-freq", rh.UpdateRebalanceFreq)

		// TODO: Update robo portfolio
		roboAdvisorGroup.PUT("/update")

		roboAdvisorGroup.DELETE("/", rh.DeleteRoboPortfolio)
	}

	manualGroup := portfolioGroup.Group("/manual-portfolio")
	{
		manualGroup.GET("/", mh.GetManualPortfolios)
		manualGroup.GET("/:name", mh.GetManualPortfolio)
		manualGroup.GET("/:name/value", mh.GetPortfolioValue)

		manualGroup.POST("/", mh.CreateManualPortfolio)
		manualGroup.POST("/:name/add", mh.AddMoneyToManualPortfolio)
		manualGroup.POST("/:name/withdraw", mh.WithDrawMoneyFromManualPortfolio)

		manualGroup.PUT("/:name", mh.UpdatePortfolioName)
		manualGroup.DELETE("/:name", mh.DeleteManualPortfolio)

		manualGroup.PUT("/:name/:symbol/buy")
		manualGroup.PUT("/:name/:symbol/sell")

		//TODO: get transactions
		manualGroup.GET("/:name/transactions")

	}
}
