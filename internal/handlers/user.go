package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/app"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var dto dto.SignUpRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := app.UserService.SignUp(dto); err != nil {
		HandleError(c, err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Sign Up Successful"})

}

func GetUser(c *gin.Context) {
	var dto dto.GetUserRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := app.UserService.GetUser(dto)
	if err != nil {
		HandleError(c, err)
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})

}
