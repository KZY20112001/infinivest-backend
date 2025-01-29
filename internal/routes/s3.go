package routes

import (
	"github.com/KZY20112001/infinivest-backend/internal/handlers"
	"github.com/KZY20112001/infinivest-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterS3Routes(r *gin.Engine, h *handlers.S3Handler) {
	s3Group := r.Group("/s3")
	s3Group.Use(middlewares.AuthMiddleware())
	{
		s3Group.POST("/upload-url", h.GeneratePresignedUploadURL)
	}

}
