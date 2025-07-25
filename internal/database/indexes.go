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

	log.Println("Database indexes created successfully")
	return nil
}
