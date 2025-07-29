package routes

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/handlers"
	"ai-analytics/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterProtectedRoutes registers protected routes that require authentication
func RegisterProtectedRoutes(r *gin.Engine, config *config.Config) {
	protectedHandler := handlers.NewProtectedHandler(config)

	protected := r.Group("/api/protected")
	protected.Use(middleware.AuthMiddleware(config))
	{
		protected.GET("/profile", protectedHandler.GetProfile)
		protected.GET("/dashboard", protectedHandler.GetDashboard)
	}
}
