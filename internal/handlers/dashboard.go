package handlers

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/models"
	"ai-analytics/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type DashboardStatsHandler struct {
	dashboardStatsService *services.DashboardStatsServcie
	config                *config.Config
}

func NewDashboardStatsHandler(db *mongo.Database, config *config.Config) *DashboardStatsHandler {
	return &DashboardStatsHandler{
		dashboardStatsService: services.NewDashboardStatsService(db, config),
		config:                config,
	}
}

func (d *DashboardStatsHandler) GetRevenueTrend(c *gin.Context) {
	var dateRange models.DateRange

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			dateRange.StartDate = startDate
		}
	}

	if endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			dateRange.EndDate = endDate
		}
	}

	revenueStats, err := d.dashboardStatsService.GetDailyRevenue(c.Request.Context(), dateRange)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"daily-revenues": revenueStats})
}
