package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "HELLO SIGN UP"})
}
