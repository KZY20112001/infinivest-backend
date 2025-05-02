package handlers

import (
	"net/http"
	"strconv"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service services.NotificationService
}

func NewNotificationHandler(ns services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: ns}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetUint("id")
	limitStr := c.Query("limit")
	limit := 0
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		limit = parsedLimit
	}
	notifications, err := h.service.GetNotifications(c.Request.Context(), userID, uint(limit))
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func (h *NotificationHandler) ClearNotifications(c *gin.Context) {
	userID := c.GetUint("id")
	err := h.service.ClearNotifications(c.Request.Context(), userID)
	if err != nil {
		commons.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notifications cleared successfully"})
}
