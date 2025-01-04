package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userHandler handlers.UserHandler, profileHandler handlers.ProfileHandler) {
	RegisterUserRoutes(r, userHandler)
	RegisterProfileRoutes(r, profileHandler)
}
