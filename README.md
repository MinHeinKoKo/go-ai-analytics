# AI Analytics for Marketing & Targeting

An advanced AI-powered analytics platform designed for Jumpstart, a nationwide fashion retailer, to enhance operational efficiency and improve customer experience through intelligent marketing and targeting solutions.

## 🎯 Project Overview

This application provides comprehensive AI analytics capabilities for marketing optimization and customer targeting, featuring:

- **Customer Segmentation**: AI-powered clustering to identify distinct customer groups
- **Behavior Prediction**: Machine learning models to predict customer churn, lifetime value, and purchase patterns
- **Campaign Optimization**: Intelligent recommendations to maximize marketing campaign performance
- **Real-time Dashboard**: Interactive analytics dashboard with key performance indicators
- **Data Management**: Comprehensive customer, purchase, and campaign data management

## 🏗️ Architecture

### Backend (Go)
- **Framework**: Gin (HTTP web framework)
- **Database**: MongoDB with optimized indexes
- **Authentication**: JWT-based authentication
- **AI/ML**: GoLearn for machine learning algorithms
- **API**: RESTful API with comprehensive endpoints

### Frontend (Next.js)
- **Framework**: Next.js 15 with TypeScript
- **UI Components**: Shadcn/ui with Radix UI primitives
- **Styling**: Tailwind CSS
- **State Management**: Zustand for global state
- **Data Fetching**: TanStack Query (React Query)
- **Charts**: Recharts for data visualization

## 🚀 Features

### AI Analytics Capabilities
1. **Customer Segmentation**
   - K-means clustering algorithm
   - Segments based on spending patterns, frequency, and demographics
   - Automatic segment profiling and insights

2. **Predictive Analytics**
   - Churn prediction with confidence scores
   - Customer lifetime value estimation
   - Next purchase timing prediction

3. **Campaign Optimization**
   - ROAS (Return on Ad Spend) optimization
   - Performance scoring and recommendations
   - Budget allocation suggestions

### Dashboard & Reporting
- Real-time metrics and KPIs
- Interactive charts and visualizations
- Customer segment analysis
- Campaign performance tracking
- Revenue and growth analytics

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

## 📋 Prerequisites

- Go 1.24.4 or higher
- Node.js 18+ and npm
- MongoDB 4.4+
- Docker (optional)

## 🛠️ Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd ai-analytics-marketing-targeting
```

### 2. Backend Setup

#### Install Go Dependencies
```bash
go mod download
```

#### Environment Configuration
Create a `.env` file in the root directory:
```env
PORT=8080
HOST=localhost
MONGODB_URI=mongodb://localhost:27017
MONGODB_DBNAME=ai-analytics
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24
```

#### Start MongoDB
```bash
# Using Docker
make docker-run

# Or start your local MongoDB instance
mongod
```

### 3. Frontend Setup

#### Navigate to web directory and install dependencies
```bash
cd web
npm install
```

#### Environment Configuration
Create a `.env.local` file in the web directory:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 🚀 Running the Application

### Start Backend Server
```bash
# From root directory
make run
# Or with live reload
make watch
```

### Start Frontend Development Server
```bash
# From web directory
cd web
npm run dev
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## 📊 Sample Data & Testing

### Generate Sample Data
1. Navigate to http://localhost:3000/login
2. Create an account or use demo credentials:
   - Email: demo@example.com
   - Password: demo123
3. Go to Analytics page and click "Generate Sample Data"
4. Explore the dashboard and analytics features

### API Testing
Use the provided test script:
```bash
chmod +x scripts/test_auth.sh
./scripts/test_auth.sh
```

## 🔧 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/me` - Get current user

### Analytics
- `GET /api/v1/analytics/dashboard` - Dashboard metrics
- `POST /api/v1/analytics/segmentation` - Customer segmentation
- `POST /api/v1/analytics/prediction` - Behavior prediction
- `POST /api/v1/analytics/optimization` - Campaign optimization

### Data Management
- `GET /api/v1/customers` - List customers
- `POST /api/v1/customers` - Create customer
- `POST /api/v1/purchases` - Create purchase
- `GET /api/v1/campaigns` - List campaigns
- `POST /api/v1/campaigns` - Create campaign
- `POST /api/v1/campaigns/performance` - Add performance data

### Utility
- `POST /api/v1/analytics/sample-data` - Generate sample data
- `POST /api/v1/analytics/import` - Import training data

## 🧠 AI/ML Implementation

### Customer Segmentation Algorithm
The application uses K-means clustering to segment customers based on:
- Total spending amount
- Purchase frequency
- Customer age
- Registration recency

### Prediction Models
1. **Churn Prediction**: Based on recency and frequency analysis
2. **Lifetime Value**: Calculated using average order value and predicted lifespan
3. **Next Purchase**: Estimated using historical purchase patterns

### Campaign Optimization
- Performance scoring based on ROAS, CTR, and conversion metrics
- Recommendation engine for budget allocation
- A/B testing suggestions for campaign improvement

## 🏢 Business Value for Jumpstart

### Operational Efficiency
- Automated customer segmentation reduces manual analysis time
- Predictive analytics enables proactive customer retention
- Campaign optimization maximizes marketing ROI

### Customer Experience Enhancement
- Personalized targeting based on AI-driven segments
- Predictive insights for better inventory management
- Improved customer journey optimization

### Data-Driven Decision Making
- Real-time analytics dashboard for quick insights
- AI-powered recommendations for marketing strategies
- Performance tracking and optimization suggestions

## 🔒 Security & Privacy

- JWT-based authentication with secure token handling
- Password hashing using bcrypt
- Input validation and sanitization
- CORS configuration for secure cross-origin requests
- Environment-based configuration management

## 🧪 Testing

### Run Backend Tests
```bash
make test
```

### Run Integration Tests
```bash
make itest
```

### Frontend Testing
```bash
cd web
npm run test
```

## 📦 Deployment

### Docker Deployment
```bash
# Build and run with Docker Compose
docker-compose up --build
```

### Production Build
```bash
# Backend
make build

# Frontend
cd web
npm run build
npm start
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- GoLearn library for machine learning capabilities
- Shadcn/ui for beautiful UI components
- TanStack Query for efficient data fetching
- Recharts for data visualization
- MongoDB for flexible data storage

---

**Note**: This is a school assignment project demonstrating AI analytics capabilities for marketing and targeting in the retail industry. The AI models are simplified for educational purposes and would require more sophisticated algorithms for production use.
