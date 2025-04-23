package routes

import (
	"github.com/gin-contrib/cors"

	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(userHandler *handlers.UserHandler, profileHandler *handlers.ProfileHandler, roboPortfolioHandler *handlers.RoboPortfolioHandler, manualPortfolioHandler *handlers.ManualPortfolioHandler, notificationHandler *handlers.NotificationHandler, s3Handler *handlers.S3Handler) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	RegisterUserRoutes(r, userHandler)
	RegisterProfileRoutes(r, profileHandler)
	RegisterPortfolioRoutes(r, roboPortfolioHandler, manualPortfolioHandler, notificationHandler)
	RegisterS3Routes(r, s3Handler)
	return r
}
