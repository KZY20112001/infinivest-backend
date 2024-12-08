package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/app"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var dto dto.AuthRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := app.UserService.SignUp(dto)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tokens": tokens})

}

func SignIn(c *gin.Context) {
	var dto dto.AuthRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := app.UserService.SignIn(dto)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

func GetUser(c *gin.Context) {
	email := c.Param("email")
	user, err := app.UserService.GetUser(email)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})

}
