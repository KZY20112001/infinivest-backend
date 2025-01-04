package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ProfileHandler interface {
	CreateProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	GetProfile(c *gin.Context)
	// GetCurrentuser(c *gin.Context)
}

type ProfileHandlerImpl struct {
	profileService services.ProfileService
}

func NewProfileHandlerImpl(ps services.ProfileService) *ProfileHandlerImpl {
	return &ProfileHandlerImpl{profileService: ps}
}

func (h *ProfileHandlerImpl) CreateProfile(c *gin.Context) {
	var dto dto.ProfileRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")

	if err := h.profileService.CreateProfile(userID, dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created the profile"})
}

func (h *ProfileHandlerImpl) UpdateProfile(c *gin.Context) {
	var dto dto.ProfileRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")

	if err := h.profileService.UpdateProfile(userID, dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated the profile"})
}

func (h *ProfileHandlerImpl) GetProfile(c *gin.Context) {
	userID := c.GetUint("id")
	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
