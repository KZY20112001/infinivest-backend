package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetCurrentUser(c *gin.Context)
}

type UserHandlerImpl struct {
	userService services.UserService
}

func NewUserHandlerImpl(us services.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{userService: us}
}

func (h *UserHandlerImpl) SignUp(c *gin.Context) {
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

func (h *UserHandlerImpl) SignIn(c *gin.Context) {
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

func (h *UserHandlerImpl) RefreshToken(c *gin.Context) {
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

func (h *UserHandlerImpl) GetCurrentUser(c *gin.Context) {
	id := c.GetUint("id")
	user, err := h.userService.GetUser(id)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})

}
