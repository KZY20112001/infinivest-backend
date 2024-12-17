package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/app"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/gin-gonic/gin"
)

func CreateProfile(c *gin.Context) {
	var dto dto.ProfileRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")

	if err := app.ProfileService.CreateProfile(userID, dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created the profile"})
}

func UpdateProfile(c *gin.Context) {
	var dto dto.ProfileRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("id")

	if err := app.ProfileService.UpdateProfile(userID, dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated the profile"})
}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("id")
	profile, err := app.ProfileService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
