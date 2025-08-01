# Technical Architecture & Code Structure

## 🏗️ System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Next.js)                      │
├─────────────────────────────────────────────────────────────────┤
│  Dashboard  │  Analytics  │  Customers  │  Campaigns  │  Auth   │
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    │ HTTP/REST API
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend (Go/Gin)                          │
├─────────────────────────────────────────────────────────────────┤
│              │              │              │              │    │
│   Handlers   │   Services   │   Models     │  Database    │Auth│
│              │              │              │              │    │
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    │ MongoDB Driver
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Database (MongoDB)                        │
├─────────────────────────────────────────────────────────────────┤
│ customers │ purchases │ campaigns │ segments │ predictions │users│
└─────────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
ai-analytics/
├── cmd/api/                    # Application entry point
│   └── main.go                # Server initialization
├── internal/                  # Private application code
│   ├── config/               # Configuration management
│   ├── database/             # Database connection & indexes
│   ├── handlers/             # HTTP request handlers
│   ├── middleware/           # HTTP middleware (auth, CORS)
│   ├── models/               # Data models & structures
│   ├── routes/               # Route definitions
│   ├── server/               # Server setup & routing
│   ├── services/             # Business logic layer
│   └── utils/                # Utility functions
├── web/                      # Frontend application
│   ├── src/
│   │   ├── app/              # Next.js app router
│   │   ├── components/       # React components
│   │   ├── hooks/            # Custom React hooks
│   │   ├── lib/              # Utilities & API client
│   │   └── store/            # State management (Zustand)
├── seed/                     # Database seeding
├── scripts/                  # Utility scripts
└── docs/                     # Documentation
```

## 🔄 Request Flow Architecture

### 1. Authentication Flow

```
Client Request ──▶ Auth Middleware ──▶ JWT Validation ──▶ User Context ──▶ Handler
     │                    │                   │               │            │
     │                    ▼                   ▼               ▼            ▼
     │              Extract Token      Verify Signature   Set User ID   Process Request
     │                    │                   │               │            │
     ▼                    ▼                   ▼               ▼            ▼
Unauthorized ◀─── Token Missing/Invalid ◀─── Expired ◀─── Success ──▶ Response
```

**Code Implementation**:
```go
// middleware/auth.go
func AuthMiddleware(config *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        if token == "" {
            c.JSON(401, gin.H{"error": "Authorization token required"})
            c.Abort()
            return
        }
        
        userID, err := utils.ValidateJWT(token, config.JWT.Secret)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", userID)
        c.Next()
    }
}
```

### 2. Analytics Request Flow

```
Frontend Request ──▶ Route Handler ──▶ Service Layer ──▶ Database ──▶ AI Procesponse
       │                   │               │              │            │              │
       ▼                   ▼               ▼              ▼            ▼              ▼
   Validation        Extract Params    Business Logic   Data Query   ML Algorithm   JSON Response
```

**Example: Customer Segmentation Flow**

```go
// handlers/analytics.go
func (h *AnalyticsHandler) PerformSegmentation(c *gin.Context) {
    var req models.SegmentationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    segments, err := h.analyticsService.PerformCustomerSegmentation(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"segments": segments})
}

// services/analytics.go
func (s *AnalyticsService) PerformCustomerSegmentation(ctx context.Context, req models.SegmentationRequest) ([]models.CustomerSegment, error) {
    customers, err := s.GetCustomers(ctx, 1000, 0)
    if err != nil {
        return nil, err
    }
    
    segments := s.performKMeansSegmentation(customers, req.Features)
    
    // Save to database
    for _, segment := range segments {
        s.db.Collection("customer_segments").InsertOne(ctx, segment)
    }
    
    return segments, nil
}
```

## 🧠 AI Analytics Service Architecture

### Service Layer Structure

```go
type AnalyticsService struct {
    db     *mongo.Database
    config *config.Config
}
```

**Core Methods**:
1. **Data Management**: CRUD operations for analytics entities
2. **AI Processing**: Machine learning algorithms implementation
3. **Aggregation**: Data summarization and metrics calculation
4. **Prediction**: Behavioral prediction algorithms

### AI Algorithm Implementation

#### 1. Customer Segmentation (K-Means)

```go
func (s *AnalyticsService) performKMeansSegmentation(customers []models.Customer, features []string) []models.CustomerSegment {
    // Feature extraction
    var totalSpents []float64
    var frequencies []int
    
    for _, customer := range customers {
        totalSpents = append(totalSpents, customer.TotalSpent)
        frequencies = append(frequencies, customer.PurchaseFrequency)
    }
    
    // Calculate thresholds (simplified K-means)
    sort.Float64s(totalSpents)
    sort.Ints(frequencies)
    
    spendThreshold1 := totalSpents[len(totalSpents)/3]
    spendThreshold2 := totalSpents[2*len(totalSpents)/3]
    
    // Segment assignment
    var highValue, mediumValue, lowValue []models.Customer
    
    for _, customer := range customers {
        score := calculateCustomerScore(customer, spendThreshold1, spendThreshold2)
        assignToSegment(customer, score, &highValue, &mediumValue, &lowValue)
    }
    
    return createSegments(highValue, mediumValue, lowValue)
}
```

#### 2. Churn Prediction

```go
func (s *AnalyticsService) predictChurn(customer models.Customer) models.PredictionResult {
    daysSinceLastPurchase := calculateDaysSinceLastPurchase(customer)
    
    var probability float64
    switch {
    case daysSinceLastPurchase > 180:
        probability = 0.8 // High risk
    case daysSinceLastPurchase > 90:
        probability = 0.5 // Medium risk
    case daysSinceLastPurchase > 30:
        probability = 0.2 // Low risk
    default:
        probability = 0.1 // Very low risk
    }
    
    // Adjust based on purchase frequency
    if customer.PurchaseFrequency > 10 {
        probability *= 0.7 // Loyal customers less likely to churn
    }
    
    return models.PredictionResult{
        CustomerID:     customer.CustomerID,
        PredictionType: "churn",
        Probability:    probability,
        Confidence:     0.75,
        CreatedAt:      time.Now(),
    }
}
```

#### 3. Campaign Optimization

```go
func (s *AnalyticsService) OptimizeCampaign(ctx context.Context, req models.CampaignOptimizationRequest) (map[string]interface{}, error) {
    // Get performance data
    performances, err := s.getCampaignPerformances(ctx, req.CampaignID)
    if err != nil {
        return nil, err
    }
    
    // Calculate metrics
    metrics := calculateCampaignMetrics(performances)
    
    // Generate recommendations based on objective
    recommendations := generateRecommendations(req.Objective, metrics)
    
    return map[string]interface{}{
        "current_metrics":    metrics,
        "recommendations":    recommendations,
        "optimization_score": calculateOptimizationScore(metrics),
    }, nil
}
```

## 🗄️ Database Architecture

### Collection Structure

```javascript
// MongoDB Collections
{
  "customers": {
    "_id": ObjectId,
    "customer_id": "CUST00001",
    "age": 25,
    "gender": "Female",
    "location": "New York",
    "income_range": "$50k-$75k",
    "registration_date": ISODate,
    "last_purchase_date": ISODate,
    "total_spent": 1250.50,
    "purchase_frequency": 8,
    "preferred_category": "Fashion"
  },
  
  "purchases": {
    "_id": ObjectId,
    "customer_id": "CUST00001",
    "product_id": "PROD001",
    "category": "Fashion",
    "amount": 89.99,
    "quantity": 1,
    "purchase_date": ISODate,
    "channel": "online"
  },
  
  "campaigns": {
    "_id": ObjectId,
    "campaign_id": "CAMP0001",
    "name": "Summer Fashion Sale",
    "type": "email",
    "target_segment": "Fashion Lovers",
    "budget": 5000.00,
    "start_date": ISODate,
    "end_date": ISODate,
    "status": "active"
  }
}
```

### Indexing Strategy

```go
// database/indexes.go
func CreateIndexes(ctx context.Context, db *mongo.Database) error {
    // Customer indexes
    customerCollection := db.Collection("customers")
    customerIndexes := []mongo.IndexModel{
        {Keys: bson.D{{Key: "customer_id", Value: 1}}, Options: options.Index().SetUnique(true)},
        {Keys: bson.D{{Key: "registration_date", Value: -1}}},
        {Keys: bson.D{{Key: "total_spent", Value: -1}}},
        {Keys: bson.D{{Key: "purchase_frequency", Value: -1}}},
    }
    
    // Purchase indexes
    purchaseCollection := db.Collection("purchases")
    purchaseIndexes := []mongo.IndexModel{
        {Keys: bson.D{{Key: "customer_id", Value: 1}}},
        {Keys: bson.D{{Key: "purchase_date", Value: -1}}},
        {Keys: bson.D{{Key: "category", Value: 1}}},
        {Keys: bson.D{{Key: "amount", Value: -1}}},
    }
    
    // Campaign indexes
    campaignCollection := db.Collection("campaigns")
    campaignIndexes := []mongo.IndexModel{
        {Keys: bson.D{{Key: "campaign_id", Value: 1}}, Options: options.Index().SetUnique(true)},
        {Keys: bson.D{{Key: "status", Value: 1}}},
        {Keys: bson.D{{Key: "start_date", Value: -1}}},
    }
    
    return createAllIndexes(ctx, customerCollection, purchaseCollection, campaignCollection)
}
```

## 🎨 Frontend Architecture

### Component Hierarchy

```
App
├── Layout
│   ├── Sidebar Navigation
│   ├── Top Bar
│   └── User Profile
├── Dashboard
│   ├── Metrics Cards
│   ├── Revenue Chart
│   ├── Customer Growth Chart
│   └── Segment Distribution
├── Analytics
│   ├── Segmentation Panel
│   ├── Prediction Panel
│   └── Optimization Panel
└── Auth
    ├── Login Form
    └── Registration Form
```

### State Management

```typescript
// store/auth.ts
interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (token: string, user: User) => void
  logout: () => void
}

// API Client with automatic token injection
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("auth_token")
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
```

### Data Fetching Strategy

```typescript
// Using TanStack Query for data management
const { data: dashboardData, isLoading } = useQuery({
  queryKey: ['dashboard'],
  queryFn: () => analyticsApi.getDashboard(),
  staleTime: 60 * 1000, // 1 minute
})

const segmentationMutation = useMutation({
  mutationFn: analyticsApi.performSegmentation,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['segments'] })
  },
})
```

## 🔐 Security Architecture

### Authentication & Authorization

```go
// JWT Token Structure
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// Token Generation
func GenerateJWT(userID, email string, secret string, expiryHours int) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiryHours))),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

### API Security Measures

1. **CORS Configuration**: Restrict cross-origin requests
2. **Rate Limiting**: Prevent API abuse
3. **Input Validation**: Sanitize all inputs
4. **SQL Injection Prevention**: Use parameterized queries
5. **Password Hashing**: bcrypt for secure password storage

## 📊 Performance Optimization

### Database Optimization

1. **Indexing Strategy**: Optimized indexes for common queries
2. **Aggregation Pipelines**: Efficient data processing
3. **Connection Pooling**: Reuse database connections
4. **Query Optimization**: Minimize database round trips

### Caching Strategy

```go
// Redis caching for predictions
func (s *AnalyticsService) getCachedPrediction(customerID string) (*models.PredictionResult, error) {
    key := fmt.Sprintf("prediction:%s", customerID)
    cached, err := s.redis.Get(key).Result()
    if err == nil {
        var prediction models.PredictionResult
        json.Unmarshal([]byte(cached), &prediction)
        return &prediction, nil
    }
    return nil, err
}
```

### Frontend Optimization

1. **Code Splitting**: Lazy load components
2. **Image Optimization**: Next.js image optimization
3. **Bundle Analysis**: Minimize JavaScript bundle size
4. **Caching**: Browser and CDN caching strategies

## 🚀 Deployment Architecture

### Production Setup

```yaml
# docker-compose.yml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - MONGODB_URI=mongodb://mongo:27017
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - mongo
      
  frontend:
    build: ./web
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://api:8080
      
  mongo:
    image: mongo:latest
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
```

### Monitoring & Logging

1. **Application Metrics**: Response times, error rates
2. **Database Monitoring**: Query performance, connection health
3. **Business Metrics**: User engagement, conversion rates
4. **Error Tracking**: Centralized error logging and alerting

This technical architecture provides a scalable, maintainable, and secure foundation for the AI Analytics platform, enabling Jumpstart to leverage advanced analytics for business growth and customer insights.
