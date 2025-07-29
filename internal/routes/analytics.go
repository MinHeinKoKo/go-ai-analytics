package routes

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/handlers"
	"ai-analytics/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAnalyticsRoutes(router *gin.Engine, db *mongo.Database, config *config.Config) {
	analyticsHandler := handlers.NewAnalyticsHandler(db, config)

	// Public routes (for demo/testing)
	public := router.Group("/api")
	{
		// Sample data generation
		public.POST("/analytics/sample-data", analyticsHandler.GenerateSampleData)
		public.POST("/analytics/import", analyticsHandler.ImportTrainingData)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(config))
	{
		// Customer management
		protected.POST("/customers", analyticsHandler.CreateCustomer)
		protected.GET("/customers", analyticsHandler.GetCustomers)

		// Purchase management
		protected.POST("/purchases", analyticsHandler.CreatePurchase)

		// Campaign management
		protected.POST("/campaigns", analyticsHandler.CreateCampaign)
		protected.GET("/campaigns", analyticsHandler.GetCampaigns)
		protected.POST("/campaigns/performance", analyticsHandler.CreateCampaignPerformance)

		// AI Analytics
		protected.POST("/analytics/segmentation", analyticsHandler.PerformSegmentation)
		protected.POST("/analytics/prediction", analyticsHandler.PredictCustomerBehavior)
		protected.POST("/analytics/optimization", analyticsHandler.OptimizeCampaign)
		protected.GET("/analytics/dashboard", analyticsHandler.GetDashboard)
	}
}
