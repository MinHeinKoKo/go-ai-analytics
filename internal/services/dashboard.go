package services

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DashboardStatsServcie struct {
	db     *mongo.Database
	config *config.Config
}

func NewDashboardStatsService(db *mongo.Database, config *config.Config) *DashboardStatsServcie {
	return &DashboardStatsServcie{
		db:     db,
		config: config,
	}
}

func (d *DashboardStatsServcie) GetDailyRevenue(ctx context.Context, dateRange models.DateRange) ([]models.DailyRevenue, error) {
	pipeline := mongo.Pipeline{
		// Match documents within the date range
		bson.D{{"$match", bson.D{
			{"purchase_date", bson.D{{"$gte", dateRange.StartDate}, {"$lte", dateRange.EndDate}}},
		}}},
		// Project a computed revenue field
		bson.D{{"$project", bson.D{
			{"day", bson.D{{"$dateTrunc", bson.D{
				{"date", "$purchase_date"},
				{"unit", "day"},
			}}}},
			{"revenue", bson.D{{"$multiply", bson.A{"$amount", "$quantity"}}}},
		}}},
		// Group by day
		bson.D{{"$group", bson.D{
			{"_id", "$day"},
			{"revenue", bson.D{{"$sum", "$revenue"}}},
		}}},
		// Sort by day ascending
		bson.D{{"$sort", bson.D{{"_id", 1}}}},
	}

	purchaseCollections := d.db.Collection("purchases")
	cursor, err := purchaseCollections.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
}
defer cursor.Close(ctx)

var results []models.DailyRevenue
if err := cursor.All(ctx, &results); err != nil {
		return nil, err
}

return results, nil
}
