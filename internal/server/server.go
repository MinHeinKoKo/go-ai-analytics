package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"

	"ai-analytics/internal/config"
	"ai-analytics/internal/database"
)

type Server struct {
	port   int
	config *config.Config
	db     *mongo.Database
}

func NewServer(config *config.Config) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	mongoDB := database.New(config)
	if mongoDB == nil {
		log.Println("Error in connection mongodb")
		return nil
	}

	NewServer := &Server{
		port:   port,
		config: config,
		db:     mongoDB,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
