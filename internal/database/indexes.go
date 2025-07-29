package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes creates necessary database indexes
func CreateIndexes(ctx context.Context, db *mongo.Database) error {
	// Create unique index on email field for users collection
	userCollection := db.Collection("users")

	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := userCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		log.Printf("Failed to create email index: %v", err)
		return err
	}

	// Create indexes for analytics collections

	// Customers collection indexes
	customerCollection := db.Collection("customers")
	customerIDIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "customer_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = customerCollection.Indexes().CreateOne(ctx, customerIDIndex)
	if err != nil {
		log.Printf("Failed to create customer_id index: %v", err)
	}

	// Purchases collection indexes
	purchaseCollection := db.Collection("purchases")
	purchaseIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "customer_id", Value: 1}}},
		{Keys: bson.D{{Key: "purchase_date", Value: -1}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
	}
	_, err = purchaseCollection.Indexes().CreateMany(ctx, purchaseIndexes)
	if err != nil {
		log.Printf("Failed to create purchase indexes: %v", err)
	}

	// Campaigns collection indexes
	campaignCollection := db.Collection("campaigns")
	campaignIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "campaign_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "start_date", Value: -1}}},
	}
	_, err = campaignCollection.Indexes().CreateMany(ctx, campaignIndexes)
	if err != nil {
		log.Printf("Failed to create campaign indexes: %v", err)
	}

	// Campaign performance collection indexes
	performanceCollection := db.Collection("campaign_performance")
	performanceIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "campaign_id", Value: 1}}},
		{Keys: bson.D{{Key: "date", Value: -1}}},
	}
	_, err = performanceCollection.Indexes().CreateMany(ctx, performanceIndexes)
	if err != nil {
		log.Printf("Failed to create performance indexes: %v", err)
	}

	// Customer segments collection indexes
	segmentCollection := db.Collection("customer_segments")
	segmentIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "segment_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = segmentCollection.Indexes().CreateOne(ctx, segmentIndex)
	if err != nil {
		log.Printf("Failed to create segment index: %v", err)
	}

	// Predictions collection indexes
	predictionCollection := db.Collection("predictions")
	predictionIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "customer_id", Value: 1}}},
		{Keys: bson.D{{Key: "prediction_type", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}
	_, err = predictionCollection.Indexes().CreateMany(ctx, predictionIndexes)
	if err != nil {
		log.Printf("Failed to create prediction indexes: %v", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
