package middlewares

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKey
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid expiration time"})
				c.Abort()
				return
			}
			if tokenType, ok := claims["type"].(string); ok {
				if tokenType != "ACCESS" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}
			if idStr, ok := claims["id"].(string); ok {
				if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
					c.Set("id", uint(id))
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}
