# UML Diagrams - AI Analytics Platform

## 1. System Architecture Diagram

```mermaid
graph TB
    subgraph "Frontend Layer"
        UI[Next.js UI Components]
        Store[Zustand State Management]
        API_Client[API Client]
    end
    
    subgraph "Backend Layer"
        Router[Gin Router]
        Auth[JWT Middleware]
        Handlers[Request Handlers]
        Services[Business Services]
        Utils[Utility Functions]
    end
    
    subgraph "Data Layer"
        MongoDB[(MongoDB Database)]
        Collections[Collections: customers, purchases, campaigns, etc.]
    end
    
    subgraph "AI/ML Layer"
        Segmentation[Customer Segmentation]
        Prediction[Predictive Analytics]
        Optimization[Campaign Optimization]
    end
    
    UI --> Store
    Store --> API_Client
    API_Client --> Router
    Router --> Auth
    Auth --> Handlers
    Handlers --> Services
    Services --> Utils
    Services --> MongoDB
    MongoDB --> Collections
    Services --> Segmentation
    Services --> Prediction
    Services --> Optimization
```

## 2. Class Diagram - Core Models

```mermaid
classDiagram
    class User {
        +ID ObjectID
        +Email string
        +Password string
        +Name string
        +CreatedAt time.Time
        +UpdatedAt time.Time
        +Validate() error
        +HashPassword() error
        +CheckPassword(password string) bool
    }
    
    class Customer {
        +ID ObjectID
        +Name string
        +Email string
        +Age int
        +Gender string
        +Location string
        +RegistrationDate time.Time
        +TotalSpent float64
        +PurchaseCount int
        +LastPurchaseDate time.Time
        +Validate() error
    }
    
    class Purchase {
        +ID ObjectID
        +CustomerID ObjectID
        +Amount float64
        +ProductCategory string
        +PurchaseDate time.Time
        +PaymentMethod string
        +Validate() error
    }
    
    class Campaign {
        +ID ObjectID
        +Name string
        +Type string
        +StartDate time.Time
        +EndDate time.Time
        +Budget float64
        +TargetAudience string
        +Status string
        +Validate() error
    }
    
    class CampaignPerformance {
        +ID ObjectID
        +CampaignID ObjectID
        +Clicks int
        +Conversions int
        +Cost float64
        +Revenue float64
        +Date time.Time
        +Validate() error
    }
    
    class CustomerSegment {
        +ID ObjectID
        +CustomerID ObjectID
        +SegmentName string
        +SegmentScore float64
        +Features map[string]float64
        +CreatedAt time.Time
        +Validate() error
    }
    
    class Prediction {
        +ID ObjectID
        +CustomerID ObjectID
        +PredictionType string
        +Value float64
        +Confidence float64
        +CreatedAt time.Time
        +ExpiresAt time.Time
        +Validate() error
    }
    
    Customer ||--o{ Purchase : "has many"
    Campaign ||--o{ CampaignPerformance : "has many"
    Customer ||--|| CustomerSegment : "belongs to"
    Customer ||--o{ Prediction : s many"
```

## 3. Service Layer Architecture

```mermaid
classDiagram
    class AnalyticsService {
        +db *mongo.Database
        +GetDashboardStats() DashboardStats
        +GetRevenueTrend() []RevenueData
        +GetCustomerGrowth() []CustomerGrowthData
        +GetCustomerSegmentDistribution() []SegmentData
        +GetCampaignPerformanceData() []CampaignData
        +GenerateCustomerSegments() error
        +PredictChurn(customerID string) ChurnPrediction
        +PredictLifetimeValue(customerID string) LTVPrediction
        +PredictNextPurchase(customerID string) NextPurchasePrediction
        +OptimizeROAS() ROASOptimization
        +MinimizeCosts() CostOptimization
        +MaximizeConversions() ConversionOptimization
    }
    
    class ImportService {
        +db *mongo.Database
        +ImportCustomers(data []Customer) ImportResult
        +ImportPurchases(data []Purchase) ImportResult
        +ImportCampaigns(data []Campaign) ImportResult
        +ImportCampaignPerformance(data []CampaignPerformance) ImportResult
        +ValidateData(data interface{}) []ValidationError
        +GenerateSampleData() SampleData
    }
    
    class AuthService {
        +db *mongo.Database
        +jwtSecret string
        +Register(user User) error
        +Login(email, password string) (string, error)
        +ValidateToken(token string) (*User, error)
        +RefreshToken(token string) (string, error)
    }
    
    class MLService {
        +PerformKMeansClustering(customers []Customer) []CustomerSegment
        +CalculateChurnProbability(customer Customer) float64
        +CalculateLifetimeValue(customer Customer) float64
        +PredictNextPurchaseDate(customer Customer) time.Time
        +AnalyzeCampaignPerformance(campaigns []Campaign) CampaignAnalysis
    }
    
    AnalyticsService --> MLService : uses
    AnalyticsService --> ImportService : collaborates
    AuthService --> User : manages
```

## 4. API Endpoints Sequence Diagram

```mermaid
sequenceDiagram
    participant Client as Frontend Client
    participant Router as Gin Router
    participant Auth as Auth Middleware
    participant Handler as Request Handler
    participant Service as Analytics Service
    participant ML as ML Service
    participant DB as MongoDB
    
    Client->>Router: GET /analytics/dashboard
    Router->>Auth: Validate JWT Token
    Auth->>Router: Token Valid
    Router->>Handler: Route to DashboardHandler
    Handler->>Service: GetDashboardStats()
    Service->>DB: Aggregate Revenue Data
    DB->>Service: Revenue Results
    Service->>DB: Count Customers
    DB->>Service: Customer Count
    Service->>ML: GetCustomerSegments()
    ML->>Service: Segment Data
    Service->>Handler: Dashboard Stats
    Handler->>Client: JSON Response
```

## 5. Data Flow Diagram

```mermaid
flowchart TD
    A[CSV File Upload] --> B[File Validation]
    B --> C{Valid Format?}
    C -->|No| D[Return Validation Errors]
    C -->|Yes| E[Parse CSV Data]
    E --> F[Data Type Validation]
    F --> G{Valid Data?}
    G -->|No| H[Return Data Errors]
    G -->|Yes| I[Transform Data]
    I --> J[Insert into MongoDB]
    J --> K[Trigger ML Processing]
    K --> L[Generate Customer Segments]
    L --> M[Calculate Predictions]
    M --> N[Update Dashboard Cache]
    N --> O[Notify Frontend]
    O --> P[Refresh Dashboard]
```

## 6. Component Interaction Diagram

```mermaid
graph LR
    subgraph "Frontend Components"
        Dashboard[Dashboard Component]
        Analytics[Analytics Component]
        DataImport[Data Import Component]
        Auth[Auth Component]
    end
    
    subgraph "State Management"
        AuthStore[Auth Store]
        DataStore[Data Store]
    end
    
    subgraph "API Layer"
        AuthAPI[Auth API]
        AnalyticsAPI[Analytics API]
        ImportAPI[Import API]
    end
    
    subgraph "Backend Services"
        AuthService[Auth Service]
        AnalyticsService[Analytics Service]
        ImportService[Import Service]
        MLService[ML Service]
    end
    
    Dashboard --> DataStore
    Analytics --> DataStore
    DataImport --> DataStore
    Auth --> AuthStore
    
    DataStore --> AnalyticsAPI
    DataStore --> ImportAPI
    AuthStore --> AuthAPI
    
    AuthAPI --> AuthService
    AnalyticsAPI --> AnalyticsService
    ImportAPI --> ImportService
    
    AnalyticsService --> MLService
```

## 7. Database Schema Diagram

```mermaid
erDiagram
    USERS {
        ObjectID _id PK
        string email UK
        string password
        string name
        datetime created_at
        datetime updated_at
    }
    
    CUSTOMERS {
        ObjectID _id PK
        string name
        string email UK
        int age
        string gender
        string location
        datetime registration_date
        float total_spent
        int purchase_count
        datetime last_purchase_date
    }
    
    PURCHASES {
        ObjectID _id PK
        ObjectID customer_id FK
        float amount
        string product_category
        datetime purchase_date
        string payment_method
    }
    
    CAMPAIGNS {
        ObjectID _id PK
        string name
        string type
        datetime start_date
        datetime end_date
        float budget
        string target_audience
        string status
    }
    
    CAMPAIGN_PERFORMANCE {
        ObjectID _id PK
        ObjectID campaign_id FK
        int clicks
        int conversions
        float cost
        float revenue
        datetime date
    }
    
    CUSTOMER_SEGMENTS {
        ObjectID _id PK
        ObjectID customer_id FK
        string segment_name
        float segment_score
        object features
        datetime created_at
    }
    
    PREDICTIONS {
        ObjectID _id PK
        ObjectID customer_id FK
        string prediction_type
        float value
        float confidence
        datetime created_at
        datetime expires_at
    }
    
    CUSTOMERS ||--o{ PURCHASES : "makes"
    CAMPAIGNS ||--o{ CAMPAIGN_PERFORMANCE : "has"
    CUSTOMERS ||--|| CUSTOMER_SEGMENTS : "belongs to"
    CUSTOMERS ||--o{ PREDICTIONS : "has"
```

## 8. Authentication Flow Diagram

```mermaid
sequenceDiagram
    participant User as User
    participant Frontend as Frontend
    participant Backend as Backend API
    participant DB as MongoDB
    participant JWT as JWT Service
    
    User->>Frontend: Enter credentials
    Frontend->>Backend: POST /auth/login
    Backend->>DB: Find user by email
    DB->>Backend: User data
    Backend->>Backend: Verify password
    Backend->>JWT: Generate JWT token
    JWT->>Backend: JWT token
    Backend->>Frontend: Return token + user data
    Frontend->>Frontend: Store token in localStorage
    Frontend->>User: Redirect to dashboard
    
    Note over Frontend,Backend: Subsequent API calls
    Frontend->>Backend: API request with JWT header
    Backend->>JWT: Validate token
    JWT->>Backend: Token valid/invalid
    Backend->>Frontend: API response or 401 error
```

## 9. ML Pipeline Architecture

```mermaid
flowchart TB
    subgraph "Data Ingestion"
        CSV[CSV Import]
        API[API Data]
        DB_Data[Existing DB Data]
    end
    
    subgraph "Data Processing"
        Clean[Data Cleaning]
        Transform[Feature Engineering]
        Validate[Data Validation]
    end
    
    subgraph "ML Models"
        Clustering[K-Means Clustering]
        Churn[Churn Prediction]
        LTV[Lifetime Value]
        NextPurchase[Next Purchase Prediction]
    end
    
    subgraph "Model Output"
        Segments[Customer Segments]
        Predictions[Predictions]
        Insights[Business Insights]
    end
    
    subgraph "Application Layer"
        Dashboard[Dashboard Updates]
        Recommendations[Recommendations]
        Alerts[Automated Alerts]
    end
    
    CSV --> Clean
    API --> Clean
    DB_Data --> Clean
    
    Clean --> Transform
    Transform --> Validate
    
    Validate --> Clustering
    Validate --> Churn
    Validate --> LTV
    Validate --> NextPurchase
    
    Clustering --> Segments
    Churn --> Predictions
    LTV --> Predictions
    NextPurchase --> Predictions
    
    Segments --> Insights
    Predictions --> Insights
    
    Insights --> Dashboard
    Insights --> Recommendations
    Insights --> Alerts
```

## 10. Deployment Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        Browser[Web Browser]
        Mobile[Mobile App]
    end
    
    subgraph "CDN/Load Balancer"
        CDN[Content Delivery Network]
        LB[Load Balancer]
    end
    
    subgraph "Application Layer"
        Frontend1[Next.js Frontend Instance 1]
        Frontend2[Next.js Frontend Instance 2]
        Backend1[Go API Instance 1]
        Backend2[Go API Instance 2]
    end
    
    subgraph "Data Layer"
        MongoDB_Primary[(MongoDB Primary)]
        MongoDB_Secondary[(MongoDB Secondary)]
        Redis[(Redis Cache)]
    end
    
    subgraph "External Services"
        Email[Email Service]
        Analytics[Analytics Service]
        Monitoring[Monitoring Service]
    end
    
    Browser --> CDN
    Mobile --> CDN
    CDN --> LB
    
    LB --> Frontend1
    LB --> Frontend2
    LB --> Backend1
    LB --> Backend2
    
    Frontend1 --> Backend1
    Frontend2 --> Backend2
    
    Backend1 --> MongoDB_Primary
    Backend2 --> MongoDB_Primary
    Backend1 --> Redis
    Backend2 --> Redis
    
    MongoDB_Primary --> MongoDB_Secondary
    
    Backend1 --> Email
    Backend2 --> Analytics
    Backend1 --> Monitoring
```

These UML diagrams provide a comprehensive view of the AI Analytics Platform architecture, covering:

1. **System Architecture**: High-level system components and their relationships
2. **Class Diagram**: Core data models and their relationships
3. **Service Layer**: Business logic organization and dependencies
4. **API Sequence**: Request/response flow through the system
5. **Data Flow**: How data moves through the import and processing pipeline
6. **Component Interaction**: Frontend component relationships and data flow
7. **Database Schema**: Entity relationships and data structure
8. **Authentication Flow**: Security and user management process
9. **ML Pipeline**: Machine learning data processing and model execution
10. **Deployment Architecture**: Production deployment structure and scaling

These diagrams serve as both documentation and architectural reference for the entire platform! 🏗️
