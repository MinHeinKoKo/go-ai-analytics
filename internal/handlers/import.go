package handlers

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/models"
	"ai-analytics/internal/services"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImportHandler struct {
	analyticsService *services.AnalyticsService
	config           *config.Config
}

func NewImportHandler(db *mongo.Database, config *config.Config) *ImportHandler {
	return &ImportHandler{
		analyticsService: services.NewAnalyticsService(db, config),
		config:           config,
	}
}

// ImportCustomers handles CSV import for customers
func (h *ImportHandler) ImportCustomers(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file type
	if !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed"})
		return
	}

	// Parse CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse CSV file"})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file must contain header and at least one data row"})
		return
	}

	// Validate headers
	expectedHeaders := []string{"customer_id", "age", "gender", "location", "income_range", "registration_date", "preferred_category"}
	headers := records[0]
	if !validateHeaders(headers, expectedHeaders) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Invalid CSV headers",
			"expected": expectedHeaders,
			"received": headers,
		})
		return
	}

	var customers []models.Customer
	var errors []string
	successCount := 0

	// Process each record
	for i, record := range records[1:] {
		customer, err := h.parseCustomerRecord(record, i+2)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: %s", i+2, err.Error()))
			continue
		}

		// Create customer
		_, err = h.analyticsService.CreateCustomer(c.Request.Context(), customer)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: Failed to save customer - %s", i+2, err.Error()))
			continue
		}

		customers = append(customers, customer)
		successCount++
	}

	response := gin.H{
		"success_count": successCount,
		"total_rows":    len(records) - 1,
		"imported":      len(customers),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}

// ImportPurchases handles CSV import for purchases
func (h *ImportHandler) ImportPurchases(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed"})
		return
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse CSV file"})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file must contain header and at least one data row"})
		return
	}

	expectedHeaders := []string{"customer_id", "product_id", "category", "amount", "quantity", "purchase_date", "channel"}
	headers := records[0]
	if !validateHeaders(headers, expectedHeaders) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Invalid CSV headers",
			"expected": expectedHeaders,
			"received": headers,
		})
		return
	}

	var purchases []models.Purchase
	var errors []string
	successCount := 0

	for i, record := range records[1:] {
		purchase, err := h.parsePurchaseRecord(record, i+2)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: %s", i+2, err.Error()))
			continue
		}

		_, err = h.analyticsService.CreatePurchase(c.Request.Context(), purchase)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: Failed to save purchase - %s", i+2, err.Error()))
			continue
		}

		purchases = append(purchases, purchase)
		successCount++
	}

	response := gin.H{
		"success_count": successCount,
		"total_rows":    len(records) - 1,
		"imported":      len(purchases),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}

// ImportCampaigns handles CSV import for campaigns
func (h *ImportHandler) ImportCampaigns(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed"})
		return
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse CSV file"})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file must contain header and at least one data row"})
		return
	}

	expectedHeaders := []string{"campaign_id", "name", "type", "target_segment", "budget", "start_date", "end_date", "status"}
	headers := records[0]
	if !validateHeaders(headers, expectedHeaders) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Invalid CSV headers",
			"expected": expectedHeaders,
			"received": headers,
		})
		return
	}

	var campaigns []models.MarketingCampaign
	var errors []string
	successCount := 0

	for i, record := range records[1:] {
		campaign, err := h.parseCampaignRecord(record, i+2)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: %s", i+2, err.Error()))
			continue
		}

		_, err = h.analyticsService.CreateCampaign(c.Request.Context(), campaign)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: Failed to save campaign - %s", i+2, err.Error()))
			continue
		}

		campaigns = append(campaigns, campaign)
		successCount++
	}

	response := gin.H{
		"success_count": successCount,
		"total_rows":    len(records) - 1,
		"imported":      len(campaigns),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}

// ImportCampaignPerformance handles CSV import for campaign performance data
func (h *ImportHandler) ImportCampaignPerformance(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed"})
		return
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse CSV file"})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file must contain header and at least one data row"})
		return
	}

	expectedHeaders := []string{"campaign_id", "impressions", "clicks", "conversions", "revenue", "cost", "date"}
	headers := records[0]
	if !validateHeaders(headers, expectedHeaders) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Invalid CSV headers",
			"expected": expectedHeaders,
			"received": headers,
		})
		return
	}

	var performances []models.CampaignPerformance
	var errors []string
	successCount := 0

	for i, record := range records[1:] {
		performance, err := h.parsePerformanceRecord(record, i+2)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: %s", i+2, err.Error()))
			continue
		}

		_, err = h.analyticsService.CreateCampaignPerformance(c.Request.Context(), performance)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: Failed to save performance data - %s", i+2, err.Error()))
			continue
		}

		performances = append(performances, performance)
		successCount++
	}

	response := gin.H{
		"success_count": successCount,
		"total_rows":    len(records) - 1,
		"imported":      len(performances),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}

// Helper functions for parsing CSV records

func (h *ImportHandler) parseCustomerRecord(record []string, rowNum int) (models.Customer, error) {
	if len(record) < 7 {
		return models.Customer{}, fmt.Errorf("insufficient columns")
	}

	age, err := strconv.Atoi(strings.TrimSpace(record[1]))
	if err != nil {
		return models.Customer{}, fmt.Errorf("invalid age: %s", record[1])
	}

	regDate, err := time.Parse("2006-01-02", strings.TrimSpace(record[5]))
	if err != nil {
		return models.Customer{}, fmt.Errorf("invalid registration date format (expected YYYY-MM-DD): %s", record[5])
	}

	customer := models.Customer{
		ID:                primitive.NewObjectID(),
		CustomerID:        strings.TrimSpace(record[0]),
		Age:               age,
		Gender:            strings.TrimSpace(record[2]),
		Location:          strings.TrimSpace(record[3]),
		IncomeRange:       strings.TrimSpace(record[4]),
		RegistrationDate:  regDate,
		PreferredCategory: strings.TrimSpace(record[6]),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return customer, nil
}

func (h *ImportHandler) parsePurchaseRecord(record []string, rowNum int) (models.Purchase, error) {
	if len(record) < 7 {
		return models.Purchase{}, fmt.Errorf("insufficient columns")
	}

	amount, err := strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
	if err != nil {
		return models.Purchase{}, fmt.Errorf("invalid amount: %s", record[3])
	}

	quantity, err := strconv.Atoi(strings.TrimSpace(record[4]))
	if err != nil {
		return models.Purchase{}, fmt.Errorf("invalid quantity: %s", record[4])
	}

	purchaseDate, err := time.Parse("2006-01-02", strings.TrimSpace(record[5]))
	if err != nil {
		return models.Purchase{}, fmt.Errorf("invalid purchase date format (expected YYYY-MM-DD): %s", record[5])
	}

	purchase := models.Purchase{
		ID:           primitive.NewObjectID(),
		CustomerID:   strings.TrimSpace(record[0]),
		ProductID:    strings.TrimSpace(record[1]),
		Category:     strings.TrimSpace(record[2]),
		Amount:       amount,
		Quantity:     quantity,
		PurchaseDate: purchaseDate,
		Channel:      strings.TrimSpace(record[6]),
		CreatedAt:    time.Now(),
	}

	return purchase, nil
}

func (h *ImportHandler) parseCampaignRecord(record []string, rowNum int) (models.MarketingCampaign, error) {
	if len(record) < 8 {
		return models.MarketingCampaign{}, fmt.Errorf("insufficient columns")
	}

	budget, err := strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
	if err != nil {
		return models.MarketingCampaign{}, fmt.Errorf("invalid budget: %s", record[4])
	}

	startDate, err := time.Parse("2006-01-02", strings.TrimSpace(record[5]))
	if err != nil {
		return models.MarketingCampaign{}, fmt.Errorf("invalid start date format (expected YYYY-MM-DD): %s", record[5])
	}

	endDate, err := time.Parse("2006-01-02", strings.TrimSpace(record[6]))
	if err != nil {
		return models.MarketingCampaign{}, fmt.Errorf("invalid end date format (expected YYYY-MM-DD): %s", record[6])
	}

	campaign := models.MarketingCampaign{
		ID:            primitive.NewObjectID(),
		CampaignID:    strings.TrimSpace(record[0]),
		Name:          strings.TrimSpace(record[1]),
		Type:          strings.TrimSpace(record[2]),
		TargetSegment: strings.TrimSpace(record[3]),
		Budget:        budget,
		StartDate:     startDate,
		EndDate:       endDate,
		Status:        strings.TrimSpace(record[7]),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return campaign, nil
}

func (h *ImportHandler) parsePerformanceRecord(record []string, rowNum int) (models.CampaignPerformance, error) {
	if len(record) < 7 {
		return models.CampaignPerformance{}, fmt.Errorf("insufficient columns")
	}

	impressions, err := strconv.Atoi(strings.TrimSpace(record[1]))
	if err != nil {
		return models.CampaignPerformance{}, fmt.Errorf("invalid impressions: %s", record[1])
	}

	clicks, err := strconv.Atoi(strings.TrimSpace(record[2]))
	if err != nil {
		return models.CampaignPerformance{}, fmt.Errorf("invalid clicks: %s", record[2])
	}

	conversions, err := strconv.Atoi(strings.TrimSpace(record[3]))
	if err != nil {
		return models.CampaignPerformance{}, fmt.Errorf("invalid conversions: %s", record[3])
	}

	revenue, err := strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
	if err != nil {
		return models.CampaignPerformance{}, fmt.Errorf("invalid revenue: %s", record[4])
	}

	cost, err := strconv.ParseFloat(strings.TrimSpace(record[5]), 64)
	if err != nil {
		return models.CampaignPerformance{}, fmt.Errorf("invalid cost: %s", record[5])
	}

	date, err := time.Parse("2006-01-02", strings.TrimSpace(record[6]))
	if err != nil {
		return models.CampaignPerformance{}, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %s", record[6])
	}

	performance := models.CampaignPerformance{
		ID:          primitive.NewObjectID(),
		CampaignID:  strings.TrimSpace(record[0]),
		Impressions: impressions,
		Clicks:      clicks,
		Conversions: conversions,
		Revenue:     revenue,
		Cost:        cost,
		Date:        date,
		CreatedAt:   time.Now(),
	}

	// Calculate metrics
	if impressions > 0 {
		performance.CTR = float64(clicks) / float64(impressions) * 100
	}
	if clicks > 0 {
		performance.CPC = cost / float64(clicks)
	}
	if cost > 0 {
		performance.ROAS = revenue / cost
	}

	return performance, nil
}

// validateHeaders checks if the CSV headers match expected headers
func validateHeaders(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}

	for i, header := range expected {
		if strings.TrimSpace(strings.ToLower(actual[i])) != strings.ToLower(header) {
			return false
		}
	}

	return true
}

// GetSampleCSV returns sample CSV content for different data types
func (h *ImportHandler) GetSampleCSV(c *gin.Context) {
	dataType := c.Param("type")

	var csvContent string
	var filename string

	switch dataType {
	case "customers":
		csvContent = `customer_id,age,gender,location,income_range,registration_date,preferred_category
CUST00001,25,Female,New York,$50k-$75k,2024-01-15,Fashion
CUST00002,35,Male,California,$75k-$100k,2024-01-20,Electronics
CUST00003,28,Female,Texas,$25k-$50k,2024-02-01,Home & Garden`
		filename = "sample_customers.csv"

	case "purchases":
		csvContent = `customer_id,product_id,category,amount,quantity,purchase_date,channel
CUST00001,PROD001,Fashion,89.99,1,2024-01-20,online
CUST00002,PROD002,Electronics,299.99,1,2024-01-25,store
CUST00003,PROD003,Home & Garden,45.50,2,2024-02-05,online`
		filename = "sample_purchases.csv"

	case "campaigns":
		csvContent = `campaign_id,name,type,target_segment,budget,start_date,end_date,status
CAMP0001,Summer Fashion Sale,email,Fashion Lovers,5000.00,2024-06-01,2024-06-30,completed
CAMP0002,Electronics Black Friday,social,Tech Enthusiasts,10000.00,2024-11-20,2024-11-30,completed
CAMP0003,Spring Collection,display,Young Adults,7500.00,2024-03-01,2024-03-31,active`
		filename = "sample_campaigns.csv"

	case "performance":
		csvContent = `campaign_id,impressions,clicks,conversions,revenue,cost,date
CAMP0001,10000,500,25,2500.00,1000.00,2024-06-01
CAMP0001,12000,600,30,3000.00,1200.00,2024-06-02
CAMP0002,15000,750,50,5000.00,1500.00,2024-11-20`
		filename = "sample_campaign_performance.csv"

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data type. Supported types: customers, purchases, campaigns, performance"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")
	c.String(http.StatusOK, csvContent)
}

// GetImportTemplate returns information about CSV import requirements
func (h *ImportHandler) GetImportTemplate(c *gin.Context) {
	templates := map[string]interface{}{
		"customers": map[string]interface{}{
			"required_headers": []string{"customer_id", "age", "gender", "location", "income_range", "registration_date", "preferred_category"},
			"data_types": map[string]string{
				"customer_id":        "string (unique identifier)",
				"age":                "integer (18-100)",
				"gender":             "string (Male/Female/Other)",
				"location":           "string (city/state)",
				"income_range":       "string (e.g., $50k-$75k)",
				"registration_date":  "date (YYYY-MM-DD format)",
				"preferred_category": "string (product category)",
			},
			"example_row": "CUST00001,25,Female,New York,$50k-$75k,2024-01-15,Fashion",
		},
		"purchases": map[string]interface{}{
			"required_headers": []string{"customer_id", "product_id", "category", "amount", "quantity", "purchase_date", "channel"},
			"data_types": map[string]string{
				"customer_id":   "string (must exist in customers)",
				"product_id":    "string (product identifier)",
				"category":      "string (product category)",
				"amount":        "decimal (purchase amount)",
				"quantity":      "integer (number of items)",
				"purchase_date": "date (YYYY-MM-DD format)",
				"channel":       "string (online/store)",
			},
			"example_row": "CUST00001,PROD001,Fashion,89.99,1,2024-01-20,online",
		},
		"campaigns": map[string]interface{}{
			"required_headers": []string{"campaign_id", "name", "type", "target_segment", "budget", "start_date", "end_date", "status"},
			"data_types": map[string]string{
				"campaign_id":    "string (unique identifier)",
				"name":           "string (campaign name)",
				"type":           "string (email/social/display/search)",
				"target_segment": "string (target audience)",
				"budget":         "decimal (campaign budget)",
				"start_date":     "date (YYYY-MM-DD format)",
				"end_date":       "date (YYYY-MM-DD format)",
				"status":         "string (active/paused/completed)",
			},
			"example_row": "CAMP0001,Summer Sale,email,Fashion Lovers,5000.00,2024-06-01,2024-06-30,completed",
		},
		"performance": map[string]interface{}{
			"required_headers": []string{"campaign_id", "impressions", "clicks", "conversions", "revenue", "cost", "date"},
			"data_types": map[string]string{
				"campaign_id": "string (must exist in campaigns)",
				"impressions": "integer (ad impressions)",
				"clicks":      "integer (ad clicks)",
				"conversions": "integer (conversions)",
				"revenue":     "decimal (revenue generated)",
				"cost":        "decimal (campaign cost)",
				"date":        "date (YYYY-MM-DD format)",
			},
			"example_row": "CAMP0001,10000,500,25,2500.00,1000.00,2024-06-01",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
		"general_guidelines": []string{
			"CSV files must include headers as the first row",
			"Date format must be YYYY-MM-DD",
			"Decimal numbers use dot (.) as separator",
			"Text fields should not contain commas",
			"File size limit: 10MB",
			"Maximum 10,000 rows per import",
		},
	})
}
