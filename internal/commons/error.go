package commons

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInternal           = errors.New("internal Error")
	ErrNil                = errors.New("nil value")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func HandleError(c *gin.Context, err error) {
	if err != nil {
		switch err {

		case ErrInternal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case ErrNil:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}
