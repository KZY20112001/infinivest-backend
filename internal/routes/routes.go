package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// user routes
	userGroup := router.Group("/user")
	{
		userGroup.POST("/signup", handlers.SignUp)
		userGroup.GET("/get/:email", handlers.GetUser)
	}

	// profile routes
	profileGroup := router.Group("/profile")
	{
		profileGroup.GET("/create")
	}

}
