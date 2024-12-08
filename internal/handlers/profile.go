package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateProfile(c *gin.Context) {
	fmt.Println("HERE IN PROFILE CREATION")
	c.JSON(http.StatusOK, gin.H{"message": "profile"})
}
