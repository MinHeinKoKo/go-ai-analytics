package database

import (
	"ai-analytics/internal/config"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
}

type service struct {
	db *mongo.Client
}

var (
	host = os.Getenv("BLUEPRINT_DB_HOST")
	port = os.Getenv("BLUEPRINT_DB_PORT")
	//database = os.Getenv("BLUEPRINT_DB_DATABASE")
)

func New(config *config.Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Create the client options
	clientOption := options.Client().ApplyURI(config.Database.URI)

	//connect to Mongo
	client, err := mongo.Connect(ctx, clientOption)

	if err != nil {
		fmt.Printf("failed to connect to MongoDB: %v", err)
		return nil
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Printf("failed to ping MongoDB: %v", err)
		return nil
	}

	db := client.Database(config.Database.Database)
	log.Printf("Connected to MongoDB database: %s", config.Database.Database)

	// Create indexes
	// if err := createMongoIndexes(ctx, db); err != nil {
	// 	fmt.Printf("failed to create MongoDB indexes: %v", err)
	// 	return nil
	// }

	return db
}
