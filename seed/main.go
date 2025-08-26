package main

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/database"
	"ai-analytics/internal/models"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	fmt.Println("🌱 Starting database seeding...")

	// Load configuration
	cfg := config.NewConfig()

	// Connect to database
	db := database.New(cfg)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	ctx := context.Background()

	// Generate sample data
	fmt.Println("📊 Generating sample data...")
	customers := generateCustomers(50)
	purchases := generatePurchases(200, customers)
	campaigns := generateMarketingCampaigns(10)
	performances := generateCampaignPerformances(campaigns)
	segments := generateCustomerSegments(5)
	predictions := generatePredictionResults(100, customers)

	// Insert data into database
	fmt.Println("💾 Inserting data into database...")

	// Insert customers
	if err := insertCustomers(ctx, db, customers); err != nil {
		log.Printf("Error inserting customers: %v", err)
	} else {
		fmt.Printf("✅ Inserted %d customers\n", len(customers))
	}

	// Insert purchases
	if err := insertPurchases(ctx, db, purchases); err != nil {
		log.Printf("Error inserting purchases: %v", err)
	} else {
		fmt.Printf("✅ Inserted %d purchases\n", len(purchases))
	}

	// Insert campaigns
	if err := insertCampaigns(ctx, db, campaigns); err != nil {
		log.Printf("Error inserting campaigns: %v", err)
	} else {
		fmt.Printf("✅ Inserted %d campaigns\n", len(campaigns))
	}

	// Insert campaign performances
	if err := insertCampaignPerformances(ctx, db, performances); err != nil {
		log.Printf("Error inserting campaign performances: %v", err)
	} else {
		fmt.Printf("✅ Inserted %d campaign performances\n", len(performances))
	}

	// Insert customer segments
	if err := insertCustomerSegments(ctx, db, segments); err != nil {
		log.Printf("Error inserting customer segments: %v", err)
	} else {
		fmt.Printf("✅ Inserted %d customer segments\n", len(segments))
	}

	// Insert prediction results
	if err := insertPredictionResults(ctx, db, predictions); err != nil {
		log.Printf("Error inserting prediction results: %v", err)
	} else {
		fmt.Printf("✅ Inserted %d prediction results\n", len(predictions))
	}

	fmt.Println("🎉 Database seeding completed successfully!")
	fmt.Println("\n📈 Summary:")
	fmt.Printf("- Customers: %d\n", len(customers))
	fmt.Printf("- Purchases: %d\n", len(purchases))
	fmt.Printf("- Campaigns: %d\n", len(campaigns))
	fmt.Printf("- Campaign Performances: %d\n", len(performances))
	fmt.Printf("- Customer Segments: %d\n", len(segments))
	fmt.Printf("- Prediction Results: %d\n", len(predictions))
}

// Database insertion functions
func insertCustomers(ctx context.Context, db *mongo.Database, customers []models.Customer) error {
	collection := db.Collection("customers")

	// Convert to interface slice for bulk insert
	docs := make([]interface{}, len(customers))
	for i, customer := range customers {
		docs[i] = customer
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

func insertPurchases(ctx context.Context, db *mongo.Database, purchases []models.Purchase) error {
	collection := db.Collection("purchases")

	docs := make([]interface{}, len(purchases))
	for i, purchase := range purchases {
		docs[i] = purchase
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

func insertCampaigns(ctx context.Context, db *mongo.Database, campaigns []models.MarketingCampaign) error {
	collection := db.Collection("campaigns")

	docs := make([]interface{}, len(campaigns))
	for i, campaign := range campaigns {
		docs[i] = campaign
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

func insertCampaignPerformances(ctx context.Context, db *mongo.Database, performances []models.CampaignPerformance) error {
	collection := db.Collection("campaign_performance")

	docs := make([]interface{}, len(performances))
	for i, performance := range performances {
		docs[i] = performance
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

func insertCustomerSegments(ctx context.Context, db *mongo.Database, segments []models.CustomerSegment) error {
	collection := db.Collection("customer_segments")

	docs := make([]interface{}, len(segments))
	for i, segment := range segments {
		docs[i] = segment
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

func insertPredictionResults(ctx context.Context, db *mongo.Database, predictions []models.PredictionResult) error {
	collection := db.Collection("predictions")

	docs := make([]interface{}, len(predictions))
	for i, prediction := range predictions {
		docs[i] = prediction
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

func generateCustomers(count int) []models.Customer {
	locations := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
	genders := []string{"Male", "Female", "Other"}
	incomeRanges := []string{"$25k-$50k", "$50k-$75k", "$75k-$100k", "$100k-$150k", "$150k+"}
	categories := []string{"Electronics", "Fashion", "Home & Garden", "Books", "Sports", "Beauty", "Automotive"}

	var customers []models.Customer
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		now := time.Now()
		regDate := now.AddDate(0, -rand.Intn(24), -rand.Intn(30)) // 0-24 months ago

		// Generate more realistic spending patterns
		age := 18 + rand.Intn(60)
		var totalSpent float64
		var purchaseFreq int

		// Age-based spending patterns
		if age < 25 {
			totalSpent = float64(50 + rand.Intn(800)) // Young adults: $50-$850
			purchaseFreq = 1 + rand.Intn(8)
		} else if age < 40 {
			totalSpent = float64(200 + rand.Intn(2000)) // Adults: $200-$2200
			purchaseFreq = 3 + rand.Intn(15)
		} else {
			totalSpent = float64(500 + rand.Intn(5000)) // Mature: $500-$5500
			purchaseFreq = 5 + rand.Intn(25)
		}

		var lastPurchase *time.Time
		if rand.Intn(10) > 1 { // 80% chance of having a purchase
			daysSinceReg := int(now.Sub(regDate).Hours() / 24)
			if daysSinceReg > 0 {
				daysAgo := rand.Intn(daysSinceReg)
				lp := now.AddDate(0, 0, -daysAgo)
				lastPurchase = &lp
			}
		}

		customers = append(customers, models.Customer{
			ID:                primitive.NewObjectID(),
			CustomerID:        fmt.Sprintf("CUST%05d", i+1),
			Age:               age,
			Gender:            genders[rand.Intn(len(genders))],
			Location:          locations[rand.Intn(len(locations))],
			IncomeRange:       incomeRanges[rand.Intn(len(incomeRanges))],
			RegistrationDate:  regDate,
			LastPurchaseDate:  lastPurchase,
			TotalSpent:        totalSpent,
			PurchaseFrequency: purchaseFreq,
			PreferredCategory: categories[rand.Intn(len(categories))],
			CreatedAt:         now,
			UpdatedAt:         now,
		})
	}
	return customers
}

func generatePurchases(count int, customers []models.Customer) []models.Purchase {
	products := []string{"PROD001", "PROD002", "PROD003", "PROD004", "PROD005", "PROD006", "PROD007", "PROD008", "PROD009", "PROD010"}
	categories := []string{"Electronics", "Fashion", "Home & Garden", "Books", "Sports", "Beauty", "Automotive"}
	channels := []string{"online", "store"}

	var purchases []models.Purchase
	now := time.Now()

	for i := 0; i < count; i++ {
		customer := customers[rand.Intn(len(customers))]

		// Generate purchase date between registration and now
		daysSinceReg := int(now.Sub(customer.RegistrationDate).Hours() / 24)
		if daysSinceReg <= 0 {
			daysSinceReg = 1
		}
		purchaseDate := customer.RegistrationDate.AddDate(0, 0, rand.Intn(daysSinceReg))

		// Category-based pricing
		category := categories[rand.Intn(len(categories))]
		var basePrice float64
		switch category {
		case "Electronics":
			basePrice = 100 + float64(rand.Intn(1000)) // $100-$1100
		case "Fashion":
			basePrice = 25 + float64(rand.Intn(200)) // $25-$225
		case "Home & Garden":
			basePrice = 50 + float64(rand.Intn(300)) // $50-$350
		case "Automotive":
			basePrice = 200 + float64(rand.Intn(800)) // $200-$1000
		default:
			basePrice = 15 + float64(rand.Intn(150)) // $15-$165
		}

		purchases = append(purchases, models.Purchase{
			ID:           primitive.NewObjectID(),
			CustomerID:   customer.CustomerID,
			ProductID:    products[rand.Intn(len(products))],
			Category:     category,
			Amount:       basePrice,
			Quantity:     1 + rand.Intn(3),
			PurchaseDate: purchaseDate,
			Channel:      channels[rand.Intn(len(channels))],
			CreatedAt:    now,
		})
	}
	return purchases
}

func generateMarketingCampaigns(count int) []models.MarketingCampaign {
	names := []string{
		"Summer Fashion Sale 2024",
		"Black Friday Electronics Blowout",
		"New Year New You Campaign",
		"Spring Home & Garden Collection",
		"Back to School Tech Deals",
		"Holiday Beauty Bonanza",
		"Winter Sports Equipment Sale",
		"Valentine's Day Special",
		"Mother's Day Gift Guide",
		"Father's Day Automotive Deals",
	}
	types := []string{"email", "social", "display", "search", "influencer"}
	segments := []string{"High Value Customers", "Young Adults 18-25", "Frequent Buyers", "At-Risk Customers", "New Customers"}
	// statuses := []string{"active", "paused", "completed"}

	var campaigns []models.MarketingCampaign
	now := time.Now()

	for i := 0; i < count; i++ {
		// Generate realistic campaign dates
		start := now.AddDate(0, -rand.Intn(12), -rand.Intn(30)) // Up to 12 months ago
		duration := 7 + rand.Intn(60)                           // 1 week to 2 months duration
		end := start.AddDate(0, 0, duration)

		// Determine status based on dates
		var status string
		if end.Before(now) {
			status = "completed"
		} else if start.Before(now) && end.After(now) {
			status = "active"
		} else {
			status = "paused"
		}

		// Budget based on campaign type
		var budget float64
		campaignType := types[rand.Intn(len(types))]
		switch campaignType {
		case "search":
			budget = 5000 + float64(rand.Intn(20000)) // $5k-$25k
		case "display":
			budget = 3000 + float64(rand.Intn(15000)) // $3k-$18k
		case "social":
			budget = 2000 + float64(rand.Intn(10000)) // $2k-$12k
		case "email":
			budget = 500 + float64(rand.Intn(3000)) // $500-$3.5k
		case "influencer":
			budget = 10000 + float64(rand.Intn(40000)) // $10k-$50k
		}

		campaigns = append(campaigns, models.MarketingCampaign{
			ID:            primitive.NewObjectID(),
			CampaignID:    fmt.Sprintf("CAMP%04d", i+1),
			Name:          names[rand.Intn(len(names))],
			Type:          campaignType,
			TargetSegment: segments[rand.Intn(len(segments))],
			Budget:        budget,
			StartDate:     start,
			EndDate:       end,
			Status:        status,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}
	return campaigns
}

func generateCampaignPerformances(campaigns []models.MarketingCampaign) []models.CampaignPerformance {
	var performances []models.CampaignPerformance
	now := time.Now()
	for _, camp := range campaigns {
		days := int(camp.EndDate.Sub(camp.StartDate).Hours() / 24)
		if days < 1 {
			days = 1
		}

		// Generate daily performance data
		for d := 0; d < days; d++ {
			date := camp.StartDate.AddDate(0, 0, d)
			impressions := 1000 + rand.Intn(10000)
			clicks := 50 + rand.Intn(impressions/10)
			conversions := 5 + rand.Intn(clicks/5)
			revenue := float64(conversions * (50 + rand.Intn(200)))
			cost := float64(impressions) * 0.001 * (0.5 + rand.Float64())

			performances = append(performances, models.CampaignPerformance{
				ID:          primitive.NewObjectID(),
				CampaignID:  camp.CampaignID,
				Impressions: impressions,
				Clicks:      clicks,
				Conversions: conversions,
				Revenue:     revenue,
				Cost:        cost,
				CTR:         float64(clicks) / float64(impressions),
				CPC:         cost / float64(clicks),
				ROAS:        revenue / cost,
				Date:        date,
				CreatedAt:   now,
			})
		}
	}
	return performances
}

func generateCustomerSegments(count int) []models.CustomerSegment {
	names := []string{
		"High Value Customers",
		"At Risk Customers",
		"New Customers",
		"Loyal Customers",
		"Discount Seekers",
	}
	descriptions := []string{
		"Customers with high lifetime value",
		"Customers likely to churn",
		"Recently acquired customers",
		"Customers with frequent purchases",
		"Customers who primarily buy discounted items",
	}

	var segments []models.CustomerSegment
	now := time.Now()
	for i := 0; i < count; i++ {
		criteria := map[string]interface{}{
			"min_purchases":        rand.Intn(5),
			"min_spend":            100 * rand.Intn(20),
			"last_active_days":     rand.Intn(90),
			"preferred_categories": []string{"Electronics", "Clothing", "Home"}[rand.Intn(3)],
		}

		segments = append(segments, models.CustomerSegment{
			ID:          primitive.NewObjectID(),
			SegmentID:   fmt.Sprintf("SEG%03d", i+1),
			Name:        names[rand.Intn(len(names))],
			Description: descriptions[rand.Intn(len(descriptions))],
			Criteria:    criteria,
			Size:        100 + rand.Intn(1000),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return segments
}

func generatePredictionResults(count int, customers []models.Customer) []models.PredictionResult {
	predictionTypes := []string{"churn", "ltv", "next_purchase"}

	var predictions []models.PredictionResult
	now := time.Now()
	for i := 0; i < count; i++ {
		customer := customers[rand.Intn(len(customers))]
		predType := predictionTypes[rand.Intn(len(predictionTypes))]

		var value float64
		switch predType {
		case "churn":
			value = rand.Float64()
		case "ltv":
			value = float64(100 + rand.Intn(5000))
		case "next_purchase":
			value = float64(1 + rand.Intn(30)) // days until next purchase
		}

		predictions = append(predictions, models.PredictionResult{
			ID:             primitive.NewObjectID(),
			CustomerID:     customer.CustomerID,
			PredictionType: predType,
			Probability:    rand.Float64(),
			Value:          value,
			Confidence:     0.7 + rand.Float64()*0.3, // 0.7-1.0
			CreatedAt:      now,
		})
	}
	return predictions
}
