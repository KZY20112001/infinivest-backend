package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// user routes
	userGroup := router.Group("/user")
	{
		userGroup.POST("/signup", handlers.SignUp)
		userGroup.POST("/signin", handlers.SignIn)
		userGroup.POST("/refresh", handlers.RefreshToken)
		userGroup.GET("/get/:email", middleware.AuthMiddleware(), handlers.GetUser)
	}

	// profile routes
	profileGroup := router.Group("/profile")
	profileGroup.Use(middleware.AuthMiddleware())
	{
		profileGroup.POST("/create", handlers.CreateProfile)
	}

}
