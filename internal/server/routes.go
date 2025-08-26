package server

import (
	"ai-analytics/internal/routes"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Register health routes
	routes.RegisterHealthRoutes(r)

	// Register authentication routes
	routes.RegisterAuthRoutes(r, s.db, s.config)

	// Register protected routes
	routes.RegisterProtectedRoutes(r, s.config)

	// Register analytics routes
	routes.RegisterAnalyticsRoutes(r, s.db, s.config)

	// Register import routes
	routes.RegisterImportRoutes(r, s.db, s.config)

	return r
}
