package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(us services.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var dto dto.AuthRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.userService.SignUp(dto)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tokens": tokens})

}

func (h *UserHandler) SignIn(c *gin.Context) {
	var dto dto.AuthRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.userService.SignIn(dto)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var dto dto.RefreshRequest
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := h.userService.RefreshRequest(dto)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	id := c.GetUint("id")
	user, err := h.userService.GetUser(id)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
