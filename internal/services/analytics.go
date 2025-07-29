package services

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/models"
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnalyticsService struct {
	db     *mongo.Database
	config *config.Config
}

func NewAnalyticsService(db *mongo.Database, config *config.Config) *AnalyticsService {
	return &AnalyticsService{
		db:     db,
		config: config,
	}
}

// Customer Analytics Methods

func (s *AnalyticsService) CreateCustomer(ctx context.Context, customer models.Customer) (*models.Customer, error) {
	customer.ID = primitive.NewObjectID()
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	collection := s.db.Collection("customers")
	_, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &customer, nil
}

func (s *AnalyticsService) GetCustomers(ctx context.Context, limit, offset int) ([]models.Customer, error) {
	collection := s.db.Collection("customers")

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}
	defer cursor.Close(ctx)

	var customers []models.Customer
	if err = cursor.All(ctx, &customers); err != nil {
		return nil, fmt.Errorf("failed to decode customers: %w", err)
	}

	return customers, nil
}

func (s *AnalyticsService) CreatePurchase(ctx context.Context, purchase models.Purchase) (*models.Purchase, error) {
	purchase.ID = primitive.NewObjectID()
	purchase.CreatedAt = time.Now()

	collection := s.db.Collection("purchases")
	_, err := collection.InsertOne(ctx, purchase)
	if err != nil {
		return nil, fmt.Errorf("failed to create purchase: %w", err)
	}

	// Update customer metrics
	go s.updateCustomerMetrics(ctx, purchase.CustomerID)

	return &purchase, nil
}

func (s *AnalyticsService) updateCustomerMetrics(ctx context.Context, customerID string) {
	// Calculate total spent and purchase frequency
	collection := s.db.Collection("purchases")

	pipeline := []bson.M{
		{"$match": bson.M{"customer_id": customerID}},
		{"$group": bson.M{
			"_id":                customerID,
			"total_spent":        bson.M{"$sum": "$amount"},
			"purchase_frequency": bson.M{"$sum": 1},
			"last_purchase_date": bson.M{"$max": "$purchase_date"},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalSpent        float64   `bson:"total_spent"`
		PurchaseFrequency int       `bson:"purchase_frequency"`
		LastPurchaseDate  time.Time `bson:"last_purchase_date"`
	}

	if cursor.Next(ctx) {
		cursor.Decode(&result)

		// Update customer record
		customerCollection := s.db.Collection("customers")
		customerCollection.UpdateOne(
			ctx,
			bson.M{"customer_id": customerID},
			bson.M{"$set": bson.M{
				"total_spent":        result.TotalSpent,
				"purchase_frequency": result.PurchaseFrequency,
				"last_purchase_date": result.LastPurchaseDate,
				"updated_at":         time.Now(),
			}},
		)
	}
}

// Campaign Analytics Methods

func (s *AnalyticsService) CreateCampaign(ctx context.Context, campaign models.MarketingCampaign) (*models.MarketingCampaign, error) {
	campaign.ID = primitive.NewObjectID()
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()

	collection := s.db.Collection("campaigns")
	_, err := collection.InsertOne(ctx, campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return &campaign, nil
}

func (s *AnalyticsService) GetCampaigns(ctx context.Context) ([]models.MarketingCampaign, error) {
	collection := s.db.Collection("campaigns")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}
	defer cursor.Close(ctx)

	var campaigns []models.MarketingCampaign
	if err = cursor.All(ctx, &campaigns); err != nil {
		return nil, fmt.Errorf("failed to decode campaigns: %w", err)
	}

	return campaigns, nil
}

func (s *AnalyticsService) CreateCampaignPerformance(ctx context.Context, performance models.CampaignPerformance) (*models.CampaignPerformance, error) {
	performance.ID = primitive.NewObjectID()
	performance.CreatedAt = time.Now()

	// Calculate metrics
	if performance.Clicks > 0 && performance.Impressions > 0 {
		performance.CTR = float64(performance.Clicks) / float64(performance.Impressions) * 100
	}
	if performance.Clicks > 0 && performance.Cost > 0 {
		performance.CPC = performance.Cost / float64(performance.Clicks)
	}
	if performance.Cost > 0 && performance.Revenue > 0 {
		performance.ROAS = performance.Revenue / performance.Cost
	}

	collection := s.db.Collection("campaign_performance")
	_, err := collection.InsertOne(ctx, performance)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign performance: %w", err)
	}

	return &performance, nil
}

// AI Analytics Methods

func (s *AnalyticsService) PerformCustomerSegmentation(ctx context.Context, req models.SegmentationRequest) ([]models.CustomerSegment, error) {
	// Get customer data
	customers, err := s.GetCustomers(ctx, 1000, 0) // Limit for demo
	if err != nil {
		return nil, fmt.Errorf("failed to get customers for segmentation: %w", err)
	}

	if len(customers) == 0 {
		return nil, errors.New("no customers found for segmentation")
	}

	// Simple K-means clustering implementation
	segments := s.performKMeansSegmentation(customers, req.Features)

	// Save segments to database
	var savedSegments []models.CustomerSegment
	for i, segment := range segments {
		segment.SegmentID = fmt.Sprintf("segment_%d", i+1)
		segment.ID = primitive.NewObjectID()
		segment.CreatedAt = time.Now()
		segment.UpdatedAt = time.Now()

		collection := s.db.Collection("customer_segments")
		_, err := collection.InsertOne(ctx, segment)
		if err != nil {
			continue
		}

		savedSegments = append(savedSegments, segment)
	}

	return savedSegments, nil
}

func (s *AnalyticsService) performKMeansSegmentation(customers []models.Customer, features []string) []models.CustomerSegment {
	// Simple 3-cluster segmentation based on spending and frequency
	// High Value: High spending, high frequency
	// Medium Value: Medium spending, medium frequency
	// Low Value: Low spending, low frequency

	var highValue, mediumValue, lowValue []models.Customer

	// Calculate thresholds
	var totalSpents []float64
	var frequencies []int

	for _, customer := range customers {
		totalSpents = append(totalSpents, customer.TotalSpent)
		frequencies = append(frequencies, customer.PurchaseFrequency)
	}

	sort.Float64s(totalSpents)
	sort.Ints(frequencies)

	spendThreshold1 := totalSpents[len(totalSpents)/3]
	spendThreshold2 := totalSpents[2*len(totalSpents)/3]
	freqThreshold1 := frequencies[len(frequencies)/3]
	freqThreshold2 := frequencies[2*len(frequencies)/3]

	for _, customer := range customers {
		score := 0
		if customer.TotalSpent > spendThreshold2 {
			score += 2
		} else if customer.TotalSpent > spendThreshold1 {
			score += 1
		}

		if customer.PurchaseFrequency > freqThreshold2 {
			score += 2
		} else if customer.PurchaseFrequency > freqThreshold1 {
			score += 1
		}

		switch {
		case score >= 3:
			highValue = append(highValue, customer)
		case score >= 1:
			mediumValue = append(mediumValue, customer)
		default:
			lowValue = append(lowValue, customer)
		}
	}

	segments := []models.CustomerSegment{
		{
			Name:        "High Value Customers",
			Description: "Customers with high spending and purchase frequency",
			Size:        len(highValue),
			Criteria: map[string]interface{}{
				"min_total_spent":        spendThreshold2,
				"min_purchase_frequency": freqThreshold2,
			},
		},
		{
			Name:        "Medium Value Customers",
			Description: "Customers with medium spending and purchase frequency",
			Size:        len(mediumValue),
			Criteria: map[string]interface{}{
				"min_total_spent":        spendThreshold1,
				"min_purchase_frequency": freqThreshold1,
			},
		},
		{
			Name:        "Low Value Customers",
			Description: "Customers with low spending and purchase frequency",
			Size:        len(lowValue),
			Criteria: map[string]interface{}{
				"max_total_spent":        spendThreshold1,
				"max_purchase_frequency": freqThreshold1,
			},
		},
	}

	return segments
}

func (s *AnalyticsService) PredictCustomerBehavior(ctx context.Context, req models.PredictionRequest) (*models.PredictionResult, error) {
	// Get customer data
	collection := s.db.Collection("customers")
	var customer models.Customer
	err := collection.FindOne(ctx, bson.M{"customer_id": req.CustomerID}).Decode(&customer)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	var prediction models.PredictionResult
	prediction.ID = primitive.NewObjectID()
	prediction.CustomerID = req.CustomerID
	prediction.PredictionType = req.PredictionType
	prediction.CreatedAt = time.Now()

	switch req.PredictionType {
	case "churn":
		prediction = s.predictChurn(customer)
	case "ltv":
		prediction = s.predictLifetimeValue(customer)
	case "next_purchase":
		prediction = s.predictNextPurchase(customer)
	default:
		return nil, errors.New("unsupported prediction type")
	}

	// Save prediction
	predictionCollection := s.db.Collection("predictions")
	_, err = predictionCollection.InsertOne(ctx, prediction)
	if err != nil {
		return nil, fmt.Errorf("failed to save prediction: %w", err)
	}

	return &prediction, nil
}

func (s *AnalyticsService) predictChurn(customer models.Customer) models.PredictionResult {
	// Simple churn prediction based on recency and frequency
	daysSinceLastPurchase := 0
	if customer.LastPurchaseDate != nil {
		daysSinceLastPurchase = int(time.Since(*customer.LastPurchaseDate).Hours() / 24)
	} else {
		daysSinceLastPurchase = 365 // No purchases
	}

	// Calculate churn probability
	var probability float64
	if daysSinceLastPurchase > 180 {
		probability = 0.8
	} else if daysSinceLastPurchase > 90 {
		probability = 0.5
	} else if daysSinceLastPurchase > 30 {
		probability = 0.2
	} else {
		probability = 0.1
	}

	// Adjust based on purchase frequency
	if customer.PurchaseFrequency > 10 {
		probability *= 0.7
	} else if customer.PurchaseFrequency < 3 {
		probability *= 1.3
	}

	if probability > 1.0 {
		probability = 1.0
	}

	return models.PredictionResult{
		CustomerID:     customer.CustomerID,
		PredictionType: "churn",
		Probability:    probability,
		Confidence:     0.75,
		CreatedAt:      time.Now(),
	}
}

func (s *AnalyticsService) predictLifetimeValue(customer models.Customer) models.PredictionResult {
	// Simple LTV prediction: average order value * purchase frequency * estimated lifespan
	avgOrderValue := customer.TotalSpent / float64(customer.PurchaseFrequency)
	if customer.PurchaseFrequency == 0 {
		avgOrderValue = 0
	}

	estimatedLifespan := 24.0 // months
	monthlyPurchaseRate := float64(customer.PurchaseFrequency) / 12.0

	ltv := avgOrderValue * monthlyPurchaseRate * estimatedLifespan

	return models.PredictionResult{
		CustomerID:     customer.CustomerID,
		PredictionType: "ltv",
		Value:          ltv,
		Confidence:     0.65,
		CreatedAt:      time.Now(),
	}
}

func (s *AnalyticsService) predictNextPurchase(customer models.Customer) models.PredictionResult {
	// Predict days until next purchase based on historical frequency
	daysBetweenPurchases := 30.0 // default
	if customer.PurchaseFrequency > 1 && customer.LastPurchaseDate != nil {
		daysSinceRegistration := time.Since(customer.RegistrationDate).Hours() / 24
		daysBetweenPurchases = daysSinceRegistration / float64(customer.PurchaseFrequency)
	}

	daysSinceLastPurchase := 0.0
	if customer.LastPurchaseDate != nil {
		daysSinceLastPurchase = time.Since(*customer.LastPurchaseDate).Hours() / 24
	}

	daysUntilNextPurchase := math.Max(0, daysBetweenPurchases-daysSinceLastPurchase)

	return models.PredictionResult{
		CustomerID:     customer.CustomerID,
		PredictionType: "next_purchase",
		Value:          daysUntilNextPurchase,
		Confidence:     0.60,
		CreatedAt:      time.Now(),
	}
}

func (s *AnalyticsService) OptimizeCampaign(ctx context.Context, req models.CampaignOptimizationRequest) (map[string]interface{}, error) {
	// Get campaign performance data
	collection := s.db.Collection("campaign_performance")
	cursor, err := collection.Find(ctx, bson.M{"campaign_id": req.CampaignID})
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign performance: %w", err)
	}
	defer cursor.Close(ctx)

	var performances []models.CampaignPerformance
	if err = cursor.All(ctx, &performances); err != nil {
		return nil, fmt.Errorf("failed to decode performance data: %w", err)
	}

	if len(performances) == 0 {
		return nil, errors.New("no performance data found for campaign")
	}

	// Calculate optimization recommendations
	recommendations := make(map[string]interface{})

	// Calculate averages
	var totalROAS, totalCTR, totalCPC float64
	var totalConversions, totalImpressions, totalClicks int
	var totalCost, totalRevenue float64

	for _, perf := range performances {
		totalROAS += perf.ROAS
		totalCTR += perf.CTR
		totalCPC += perf.CPC
		totalConversions += perf.Conversions
		totalImpressions += perf.Impressions
		totalClicks += perf.Clicks
		totalCost += perf.Cost
		totalRevenue += perf.Revenue
	}

	count := float64(len(performances))
	avgROAS := totalROAS / count
	avgCTR := totalCTR / count
	avgCPC := totalCPC / count

	recommendations["current_metrics"] = map[string]interface{}{
		"avg_roas":          avgROAS,
		"avg_ctr":           avgCTR,
		"avg_cpc":           avgCPC,
		"total_conversions": totalConversions,
		"total_revenue":     totalRevenue,
		"total_cost":        totalCost,
	}

	// Generate recommendations based on objective
	switch req.Objective {
	case "maximize_roas":
		recommendations["recommendations"] = []string{
			"Focus budget on high-performing segments",
			"Reduce spend on low ROAS keywords/audiences",
			"Increase bids for high-converting demographics",
		}
		if avgROAS < 2.0 {
			recommendations["priority_actions"] = []string{
				"Review and optimize targeting criteria",
				"Improve ad creative and messaging",
				"Consider pausing underperforming ad sets",
			}
		}
	case "minimize_cost":
		recommendations["recommendations"] = []string{
			"Lower bids on expensive keywords",
			"Focus on organic reach opportunities",
			"Optimize ad scheduling for peak performance hours",
		}
		recommendations["suggested_budget_reduction"] = totalCost * 0.15 // 15% reduction
	case "maximize_conversions":
		recommendations["recommendations"] = []string{
			"Increase budget for high-converting campaigns",
			"Expand successful audience segments",
			"Test new ad formats and placements",
		}
		if totalConversions < 100 {
			recommendations["priority_actions"] = []string{
				"Review conversion tracking setup",
				"Optimize landing page experience",
				"Test different call-to-action messages",
			}
		}
	}

	recommendations["optimization_score"] = s.calculateOptimizationScore(avgROAS, avgCTR, totalConversions)

	return recommendations, nil
}

func (s *AnalyticsService) calculateOptimizationScore(roas, ctr float64, conversions int) float64 {
	// Simple scoring algorithm (0-100)
	score := 0.0

	// ROAS component (40% weight)
	if roas >= 4.0 {
		score += 40
	} else if roas >= 2.0 {
		score += 30
	} else if roas >= 1.0 {
		score += 20
	} else {
		score += 10
	}

	// CTR component (30% weight)
	if ctr >= 3.0 {
		score += 30
	} else if ctr >= 2.0 {
		score += 25
	} else if ctr >= 1.0 {
		score += 20
	} else {
		score += 10
	}

	// Conversions component (30% weight)
	if conversions >= 100 {
		score += 30
	} else if conversions >= 50 {
		score += 25
	} else if conversions >= 10 {
		score += 20
	} else {
		score += 10
	}

	return score
}

func (s *AnalyticsService) GetAnalyticsDashboard(ctx context.Context, dateRange models.DateRange) (map[string]interface{}, error) {
	dashboard := make(map[string]interface{})

	// Customer metrics
	customerCollection := s.db.Collection("customers")
	totalCustomers, _ := customerCollection.CountDocuments(ctx, bson.M{})

	// Purchase metrics
	purchaseCollection := s.db.Collection("purchases")
	purchaseFilter := bson.M{}
	if !dateRange.StartDate.IsZero() && !dateRange.EndDate.IsZero() {
		purchaseFilter["purchase_date"] = bson.M{
			"$gte": dateRange.StartDate,
			"$lte": dateRange.EndDate,
		}
	}

	totalPurchases, _ := purchaseCollection.CountDocuments(ctx, purchaseFilter)

	// Revenue calculation
	pipeline := []bson.M{
		{"$match": purchaseFilter},
		{"$group": bson.M{
			"_id":           nil,
			"total_revenue": bson.M{"$sum": "$amount"},
			"avg_order":     bson.M{"$avg": "$amount"},
		}},
	}

	cursor, err := purchaseCollection.Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		var result struct {
			TotalRevenue float64 `bson:"total_revenue"`
			AvgOrder     float64 `bson:"avg_order"`
		}
		if cursor.Next(ctx) {
			cursor.Decode(&result)
			dashboard["total_revenue"] = result.TotalRevenue
			dashboard["avg_order_value"] = result.AvgOrder
		}
	}

	// Campaign metrics
	campaignCollection := s.db.Collection("campaigns")
	totalCampaigns, _ := campaignCollection.CountDocuments(ctx, bson.M{})
	activeCampaigns, _ := campaignCollection.CountDocuments(ctx, bson.M{"status": "active"})

	dashboard["total_customers"] = totalCustomers
	dashboard["total_purchases"] = totalPurchases
	dashboard["total_campaigns"] = totalCampaigns
	dashboard["active_campaigns"] = activeCampaigns

	return dashboard, nil
}
