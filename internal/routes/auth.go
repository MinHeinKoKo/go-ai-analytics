package routes

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/handlers"
	"ai-analytics/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterAuthRoutes registers all authentication routes
func RegisterAuthRoutes(r *gin.Engine, db *mongo.Database, config *config.Config) {
	authHandler := handlers.NewAuthHandler(db, config)

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.AuthMiddleware(config), authHandler.GetMe)
	}
}
