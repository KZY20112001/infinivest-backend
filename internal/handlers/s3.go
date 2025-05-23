package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type S3Handler struct {
	s3Service services.S3Service
}

func NewS3Handler(s3Service services.S3Service) *S3Handler {
	return &S3Handler{s3Service: s3Service}
}

func (h *S3Handler) GeneratePresignedUploadURL(c *gin.Context) {
	var req dto.PresignedUploadUrlRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	presignedUrl, err := h.s3Service.GeneratePresignedUploadURL(c, req)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": presignedUrl})

}
