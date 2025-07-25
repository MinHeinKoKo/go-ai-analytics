package config

import (
	"ai-analytics/internal/helpers"
	"net/http"
	"strings"
)

type Config struct {
	// Add your config fields here
	// You can use the helper functions from the helpers package
	Port     int          `json:"port"`
	Host     string       `json:"host"`
	Database MongoDbCofig `json:"database"`
	Kafka    KafkaConfig  `json:"kafka"`
	JWT      JWTConfig    `json:"jwt"`
}

type MongoDbCofig struct {
	URI      string `json:"uri"`
	Database string `json:"database"`
}

type KafkaConfig struct {
	Brokers             []string `json:"brokers"`
	TopicCrawledContent string   `json:"topic_crawled_content"`
	TopicBlogEvents     string   `json:"topic_blog_events"`
	ConsumerGroup       string   `json:"consumer_group"`
}

type JWTConfig struct {
	Secret      string `json:"secret"`
	ExpiryHours int    `json:"expiry_hours"`
}

// NewConfig creates a new config instance with values from environment variables
func NewConfig() *Config {
	return &Config{
		Port: helpers.GetEnvAsInt("PORT", 8080),
		Host: helpers.GetEnv("HOST", "localhost"),
		Database: MongoDbCofig{
			URI:      helpers.GetEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: helpers.GetEnv("MONGODB_DBNAME", "ai-analytics"),
		},
		Kafka: KafkaConfig{
			Brokers:             strings.Split(helpers.GetEnv("KAFKA_BROKERS", "localhost:9092"), ","),
			TopicCrawledContent: helpers.GetEnv("KAFKA_TOPIC_CRAWLED_CONTENT", "crawled-content"),
			TopicBlogEvents:     helpers.GetEnv("KAFKA_TOPIC_BLOG_EVENTS", "blog-events"),
			ConsumerGroup:       helpers.GetEnv("KAFKA_CONSUMER_GROUP", "blog-service"),
		},
		JWT: JWTConfig{
			Secret:      helpers.GetEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			ExpiryHours: helpers.GetEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
	}
}

// JSON helper methods that use the helpers package
func (c *Config) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	return helpers.ReadJSON(w, r, data)
}

func (c *Config) WriteJSON(w http.ResponseWriter, data any, status int, headers ...http.Header) error {
	return helpers.WriteJSON(w, data, status, headers...)
}

func (c *Config) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	return helpers.ErrorJSON(w, err, status...)
}
