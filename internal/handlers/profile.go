package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "profile"})

}
