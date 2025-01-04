package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterProfileRoutes(r *gin.Engine, h handlers.ProfileHandler) {
	profileGroup := r.Group("/profile")
	profileGroup.Use(middlewares.AuthMiddleware())
	{
		profileGroup.POST("/", h.CreateProfile)
		profileGroup.PATCH("/", h.UpdateProfile)
		profileGroup.GET("/", h.GetProfile)
	}
}
