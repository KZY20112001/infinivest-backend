package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, h handlers.UserHandler) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/signup", h.SignUp)
		userGroup.POST("/signin", h.SignIn)
		userGroup.POST("/refresh", h.RefreshToken)
		userGroup.GET("", middlewares.AuthMiddleware(), h.GetCurrentUser)
	}

}
