package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService services.ProfileService
}

func NewProfileHandler(ps services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: ps}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var dto dto.ProfileRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")

	if err := h.profileService.CreateProfile(userID, dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created the profile"})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var dto dto.ProfileRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")

	if err := h.profileService.UpdateProfile(userID, dto); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated the profile"})
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("id")
	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
