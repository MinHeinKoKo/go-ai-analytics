package routes

import (
	"ai-analytics/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterHealthRoutes registers health check routes
func RegisterHealthRoutes(r *gin.Engine) {
	healthHandler := handlers.NewHealthHandler()

	r.GET("/health", healthHandler.GetHealth)
}
