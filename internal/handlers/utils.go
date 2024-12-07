package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/global"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	if err != nil {
		switch err {
		case global.ErrDuplicate:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case global.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case global.ErrInternal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case global.ErrNil:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}
