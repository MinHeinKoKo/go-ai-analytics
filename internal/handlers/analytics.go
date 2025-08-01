package handlers

import (
	"ai-analytics/internal/config"
	"ai-analytics/internal/models"
	"ai-analytics/internal/services"
	"ai-analytics/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	config           *config.Config
}

func NewAnalyticsHandler(db *mongo.Database, config *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: services.NewAnalyticsService(db, config),
		config:           config,
	}
}

// Customer Management

func (h *AnalyticsHandler) CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCustomer, err := h.analyticsService.CreateCustomer(c.Request.Context(), customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"customer": createdCustomer})
}

func (h *AnalyticsHandler) GetCustomers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	customers, err := h.analyticsService.GetCustomers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"customers": customers})
}

func (h *AnalyticsHandler) CreatePurchase(c *gin.Context) {
	var purchase models.Purchase
	if err := c.ShouldBindJSON(&purchase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(purchase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPurchase, err := h.analyticsService.CreatePurchase(c.Request.Context(), purchase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"purchase": createdPurchase})
}

// Campaign Management

func (h *AnalyticsHandler) CreateCampaign(c *gin.Context) {
	var campaign models.MarketingCampaign
	if err := c.ShouldBindJSON(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCampaign, err := h.analyticsService.CreateCampaign(c.Request.Context(), campaign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"campaign": createdCampaign})
}

func (h *AnalyticsHandler) GetCampaigns(c *gin.Context) {
	campaigns, err := h.analyticsService.GetCampaigns(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}

func (h *AnalyticsHandler) CreateCampaignPerformance(c *gin.Context) {
	var performance models.CampaignPerformance
	if err := c.ShouldBindJSON(&performance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(performance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPerformance, err := h.analyticsService.CreateCampaignPerformance(c.Request.Context(), performance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"performance": createdPerformance})
}

// AI Analytics

func (h *AnalyticsHandler) PerformSegmentation(c *gin.Context) {
	var req models.SegmentationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	segments, err := h.analyticsService.PerformCustomerSegmentation(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"segments": segments})
}

func (h *AnalyticsHandler) PredictCustomerBehavior(c *gin.Context) {
	var req models.PredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prediction, err := h.analyticsService.PredictCustomerBehavior(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"prediction": prediction})
}

func (h *AnalyticsHandler) OptimizeCampaign(c *gin.Context) {
	var req models.CampaignOptimizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	optimization, err := h.analyticsService.OptimizeCampaign(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"optimization": optimization})
}

func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	var dateRange models.DateRange

	// Parse optional date range parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			dateRange.StartDate = startDate
		}
	}

	if endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			dateRange.EndDate = endDate
		}
	}

	dashboard, err := h.analyticsService.GetAnalyticsDashboard(c.Request.Context(), dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dashboard": dashboard})
}

// Bulk Data Import for Training

func (h *AnalyticsHandler) ImportTrainingData(c *gin.Context) {
	var data struct {
		Customers   []models.Customer            `json:"customers"`
		Purchases   []models.Purchase            `json:"purchases"`
		Campaigns   []models.MarketingCampaign   `json:"campaigns"`
		Performance []models.CampaignPerformance `json:"performance"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results := make(map[string]interface{})

	// Import customers
	if len(data.Customers) > 0 {
		var imported int
		for _, customer := range data.Customers {
			if _, err := h.analyticsService.CreateCustomer(c.Request.Context(), customer); err == nil {
				imported++
			}
		}
		results["customers_imported"] = imported
	}

	// Import purchases
	if len(data.Purchases) > 0 {
		var imported int
		for _, purchase := range data.Purchases {
			if _, err := h.analyticsService.CreatePurchase(c.Request.Context(), purchase); err == nil {
				imported++
			}
		}
		results["purchases_imported"] = imported
	}

	// Import campaigns
	if len(data.Campaigns) > 0 {
		var imported int
		for _, campaign := range data.Campaigns {
			if _, err := h.analyticsService.CreateCampaign(c.Request.Context(), campaign); err == nil {
				imported++
			}
		}
		results["campaigns_imported"] = imported
	}

	// Import performance data
	if len(data.Performance) > 0 {
		var imported int
		for _, perf := range data.Performance {
			if _, err := h.analyticsService.CreateCampaignPerformance(c.Request.Context(), perf); err == nil {
				imported++
			}
		}
		results["performance_imported"] = imported
	}

	c.JSON(http.StatusOK, gin.H{"import_results": results})
}

// Generate Sample Data for Testing

func (h *AnalyticsHandler) GenerateSampleData(c *gin.Context) {
	// Generate sample customers
	sampleCustomers := []models.Customer{
		{
			CustomerID:        "CUST001",
			Age:               25,
			Gender:            "Female",
			Location:          "New York",
			IncomeRange:       "$50k-$75k",
			RegistrationDate:  time.Now().AddDate(0, -6, 0),
			TotalSpent:        1250.50,
			PurchaseFrequency: 8,
			PreferredCategory: "Fashion",
		},
		{
			CustomerID:        "CUST002",
			Age:               35,
			Gender:            "Male",
			Location:          "California",
			IncomeRange:       "$75k-$100k",
			RegistrationDate:  time.Now().AddDate(0, -12, 0),
			TotalSpent:        2100.75,
			PurchaseFrequency: 15,
			PreferredCategory: "Electronics",
		},
		{
			CustomerID:        "CUST003",
			Age:               28,
			Gender:            "Female",
			Location:          "Texas",
			IncomeRange:       "$25k-$50k",
			RegistrationDate:  time.Now().AddDate(0, -3, 0),
			TotalSpent:        450.25,
			PurchaseFrequency: 3,
			PreferredCategory: "Home & Garden",
		},
	}

	// Generate sample purchases
	samplePurchases := []models.Purchase{
		{
			CustomerID:   "CUST001",
			ProductID:    "PROD001",
			Category:     "Fashion",
			Amount:       89.99,
			Quantity:     1,
			PurchaseDate: time.Now().AddDate(0, 0, -5),
			Channel:      "online",
		},
		{
			CustomerID:   "CUST002",
			ProductID:    "PROD002",
			Category:     "Electronics",
			Amount:       299.99,
			Quantity:     1,
			PurchaseDate: time.Now().AddDate(0, 0, -10),
			Channel:      "store",
		},
	}

	// Generate sample campaigns
	sampleCampaigns := []models.MarketingCampaign{
		{
			CampaignID:    "CAMP001",
			Name:          "Summer Fashion Sale",
			Type:          "email",
			TargetSegment: "Fashion Lovers",
			Budget:        5000.00,
			StartDate:     time.Now().AddDate(0, 0, -30),
			EndDate:       time.Now().AddDate(0, 0, -15),
			Status:        "completed",
		},
		{
			CampaignID:    "CAMP002",
			Name:          "Electronics Black Friday",
			Type:          "social",
			TargetSegment: "Tech Enthusiasts",
			Budget:        10000.00,
			StartDate:     time.Now().AddDate(0, 0, -7),
			EndDate:       time.Now().AddDate(0, 0, 7),
			Status:        "active",
		},
	}

	results := make(map[string]interface{})

	// Create sample data
	for _, customer := range sampleCustomers {
		h.analyticsService.CreateCustomer(c.Request.Context(), customer)
	}
	results["customers_created"] = len(sampleCustomers)

	for _, purchase := range samplePurchases {
		h.analyticsService.CreatePurchase(c.Request.Context(), purchase)
	}
	results["purchases_created"] = len(samplePurchases)

	for _, campaign := range sampleCampaigns {
		h.analyticsService.CreateCampaign(c.Request.Context(), campaign)
	}
	results["campaigns_created"] = len(sampleCampaigns)

	c.JSON(http.StatusOK, gin.H{"sample_data": results})
}

// Advanced Customer Lifetime Value Prediction
func (h *AnalyticsHandler) PredictLifetimeValue(c *gin.Context) {
	customerID := c.Param("customerID")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	prediction, err := h.analyticsService.PredictLifetimeValueAdvanced(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"prediction": prediction})
}

// Advanced Next Purchase Prediction
func (h *AnalyticsHandler) PredictNextPurchase(c *gin.Context) {
	customerID := c.Param("customerID")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	prediction, err := h.analyticsService.PredictNextPurchaseAdvanced(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"prediction": prediction})
}

// Campaign Cost Minimization
func (h *AnalyticsHandler) MinimizeCampaignCost(c *gin.Context) {
	campaignID := c.Param("campaignID")
	if campaignID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign ID is required"})
		return
	}

	optimization, err := h.analyticsService.MinimizeCampaignCost(c.Request.Context(), campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"optimization": optimization})
}

// Campaign Conversion Maximization
func (h *AnalyticsHandler) MaximizeCampaignConversions(c *gin.Context) {
	campaignID := c.Param("campaignID")
	if campaignID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign ID is required"})
		return
	}

	optimization, err := h.analyticsService.MaximizeCampaignConversions(c.Request.Context(), campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"optimization": optimization})
}
