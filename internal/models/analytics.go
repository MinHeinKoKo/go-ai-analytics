package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer represents customer data for analytics
type Customer struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID        string             `json:"customer_id" bson:"customer_id"`
	Age               int                `json:"age" bson:"age"`
	Gender            string             `json:"gender" bson:"gender"`
	Location          string             `json:"location" bson:"location"`
	IncomeRange       string             `json:"income_range" bson:"income_range"`
	RegistrationDate  time.Time          `json:"registration_date" bson:"registration_date"`
	LastPurchaseDate  *time.Time         `json:"last_purchase_date" bson:"last_purchase_date"`
	TotalSpent        float64            `json:"total_spent" bson:"total_spent"`
	PurchaseFrequency int                `json:"purchase_frequency" bson:"purchase_frequency"`
	PreferredCategory string             `json:"preferred_category" bson:"preferred_category"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

// Purchase represents purchase transaction data
type Purchase struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID   string             `json:"customer_id" bson:"customer_id"`
	ProductID    string             `json:"product_id" bson:"product_id"`
	Category     string             `json:"category" bson:"category"`
	Amount       float64            `json:"amount" bson:"amount"`
	Quantity     int                `json:"quantity" bson:"quantity"`
	PurchaseDate time.Time          `json:"purchase_date" bson:"purchase_date"`
	Channel      string             `json:"channel" bson:"channel"` // online, store
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

type DailyRevenue struct {
	Date    time.Time `json:"date" bson:"_id"` // _id will be the grouped date
	Revenue float64   `json:"revenue" bson:"revenue"`
}


// MarketingCampaign represents marketing campaign data
type MarketingCampaign struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CampaignID    string             `json:"campaign_id" bson:"campaign_id"`
	Name          string             `json:"name" bson:"name"`
	Type          string             `json:"type" bson:"type"` // email, social, display, search
	TargetSegment string             `json:"target_segment" bson:"target_segment"`
	Budget        float64            `json:"budget" bson:"budget"`
	StartDate     time.Time          `json:"start_date" bson:"start_date"`
	EndDate       time.Time          `json:"end_date" bson:"end_date"`
	Status        string             `json:"status" bson:"status"` // active, paused, completed
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

// CampaignPerformance represents campaign performance metrics
type CampaignPerformance struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CampaignID  string             `json:"campaign_id" bson:"campaign_id"`
	Impressions int                `json:"impressions" bson:"impressions"`
	Clicks      int                `json:"clicks" bson:"clicks"`
	Conversions int                `json:"conversions" bson:"conversions"`
	Revenue     float64            `json:"revenue" bson:"revenue"`
	Cost        float64            `json:"cost" bson:"cost"`
	CTR         float64            `json:"ctr" bson:"ctr"`   // Click-through rate
	CPC         float64            `json:"cpc" bson:"cpc"`   // Cost per click
	ROAS        float64            `json:"roas" bson:"roas"` // Return on ad spend
	Date        time.Time          `json:"date" bson:"date"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

// CustomerSegment represents AI-generated customer segments
type CustomerSegment struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	SegmentID   string                 `json:"segment_id" bson:"segment_id"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Criteria    map[string]interface{} `json:"criteria" bson:"criteria"`
	Size        int                    `json:"size" bson:"size"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// PredictionResult represents AI prediction results
type PredictionResult struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID     string             `json:"customer_id" bson:"customer_id"`
	PredictionType string             `json:"prediction_type" bson:"prediction_type"` // churn, ltv, next_purchase
	Probability    float64            `json:"probability" bson:"probability"`
	Value          float64            `json:"value" bson:"value"`
	Confidence     float64            `json:"confidence" bson:"confidence"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

// AnalyticsRequest represents request for analytics operations
type AnalyticsRequest struct {
	Type       string                 `json:"type" validate:"required"`
	Parameters map[string]interface{} `json:"parameters"`
	DateRange  DateRange              `json:"date_range"`
}

// DateRange represents date range for analytics
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// SegmentationRequest represents customer segmentation request
type SegmentationRequest struct {
	Algorithm  string                 `json:"algorithm" validate:"required"` // kmeans, dbscan
	Features   []string               `json:"features" validate:"required"`
	Parameters map[string]interface{} `json:"parameters"`
}

// PredictionRequest represents prediction request
type PredictionRequest struct {
	CustomerID     string `json:"customer_id" validate:"required"`
	PredictionType string `json:"prediction_type" validate:"required"`
}

// CampaignOptimizationRequest represents campaign optimization request
type CampaignOptimizationRequest struct {
	CampaignID string                 `json:"campaign_id" validate:"required"`
	Objective  string                 `json:"objective" validate:"required"` // maximize_roas, minimize_cost, maximize_conversions
	Parameters map[string]interface{} `json:"parameters"`
}
