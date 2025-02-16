package handlers

import (
	"net/http"

	"github.com/KZY20112001/infinivest-backend/internal/constants"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	if err != nil {
		switch err {

		case constants.ErrInternal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case constants.ErrNil:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}
