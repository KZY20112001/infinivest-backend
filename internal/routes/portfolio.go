package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterPortfolioRoutes(r *gin.Engine, rh *handlers.RoboPortfolioHandler, mh *handlers.ManualPortfolioHandler, nh *handlers.NotificationHandler) {
	portfolioGroup := r.Group("/portfolio")
	portfolioGroup.Use(middlewares.AuthMiddleware())
	notificationGroup := portfolioGroup.Group("/notifications")
	{
		notificationGroup.GET("/", nh.GetNotifications)
		notificationGroup.DELETE("/", nh.ClearNotifications)
	}

	roboAdvisorGroup := portfolioGroup.Group("/robo-portfolio")
	{
		roboAdvisorGroup.GET("/details", rh.GetRoboPortfolioDetails)
		roboAdvisorGroup.GET("/summary", rh.GetRoboPortfolioSummary)

		roboAdvisorGroup.POST("/generate/categories", rh.GenerateRoboAdvisorPortfolio)
		roboAdvisorGroup.POST("/generate/assets", rh.GenerateAssetAllocation)
		roboAdvisorGroup.POST("/confirm", rh.ConfirmGeneratedRoboPortfolio)
		roboAdvisorGroup.POST("/add", rh.AddMoneyToRoboPortfolio)
		roboAdvisorGroup.POST("/withdraw", rh.WithDrawMoneyFromRoboPortfolio)
		roboAdvisorGroup.PUT("/rebalance-freq", rh.UpdateRebalanceFreq)

		// TODO: Update robo portfolio
		roboAdvisorGroup.PUT("/update")

		roboAdvisorGroup.DELETE("/", rh.DeleteRoboPortfolio)

		roboAdvisorGroup.GET("/transactions", rh.GetRoboPortfolioTransactions)
		roboAdvisorGroup.GET("/rebalance/details", rh.GetRebalanceEvents)
		roboAdvisorGroup.PATCH("/rebalance/seen", rh.UpdateLastSeenRebalanceEvent)

		// testing only
		roboAdvisorGroup.GET("/rebalance", rh.RebalanceRoboPortfolio)
	}

	manualGroup := portfolioGroup.Group("/manual-portfolio")
	{
		manualGroup.GET("/details", mh.GetManualPortfoliosDetails)
		manualGroup.GET("/summaries", mh.GetManualPortfoliosSummaries)
		manualGroup.GET("/:name", mh.GetManualPortfolio)
		manualGroup.GET("/:name/value", mh.GetPortfolioValue)

		manualGroup.POST("/", mh.CreateManualPortfolio)
		manualGroup.POST("/:name/add", mh.AddMoneyToManualPortfolio)
		manualGroup.POST("/:name/withdraw", mh.WithDrawMoneyFromManualPortfolio)

		manualGroup.PUT("/:name", mh.UpdatePortfolioName)
		manualGroup.DELETE("/:name", mh.DeleteManualPortfolio)

		manualGroup.PUT("/:name/buy", mh.BuyAssetForManualPortfolio)
		manualGroup.PUT("/:name/sell", mh.SellAssetForManualPortfolio)

		manualGroup.GET("/:name/transactions", mh.GetManualPortfolioTransactions)
	}
}
