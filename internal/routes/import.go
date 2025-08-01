package routes

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/handlers"
	"ai-analytics/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterImportRoutes(router *gin.Engine, db *mongo.Database, config *config.Config) {
	importHandler := handlers.NewImportHandler(db, config)

	// Public routes for templates and samples
	public := router.Group("/api/v1/import")
	{
		// Get import templates and requirements
		public.GET("/templates", importHandler.GetImportTemplate)

		// Download sample CSV files
		public.GET("/sample/:type", importHandler.GetSampleCSV)
	}

	// Protected routes for actual imports
	protected := router.Group("/api/v1/import")
	protected.Use(middleware.AuthMiddleware(config))
	{
		// CSV import endpoints
		protected.POST("/customers", importHandler.ImportCustomers)
		protected.POST("/purchases", importHandler.ImportPurchases)
		protected.POST("/campaigns", importHandler.ImportCampaigns)
		protected.POST("/performance", importHandler.ImportCampaignPerformance)
	}
}
