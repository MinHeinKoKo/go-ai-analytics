package handlers

import (
	"ai-analytics/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProtectedHandler struct {
	config *config.Config
}

func NewProtectedHandler(config *config.Config) *ProtectedHandler {
	return &ProtectedHandler{
		config: config,
	}
}

// GetProfile returns user profile information
func (h *ProtectedHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	userEmail := c.MustGet("user_email").(string)

	c.JSON(http.StatusOK, gin.H{
		"message":    "This is a protected route",
		"user_id":    userID,
		"user_email": userEmail,
	})
}

// GetDashboard returns dashboard data
func (h *ProtectedHandler) GetDashboard(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your dashboard",
		"user_id": userID,
		"data": gin.H{
			"stats": gin.H{
				"total_posts": 42,
				"total_views": 1337,
			},
		},
	})
}
