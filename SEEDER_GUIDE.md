# Database Seeder Guide

## Overview
The database seeder (`seed/main.go`) generates realistic sample data for the AI Analytics platform and inserts it directly into MongoDB. This is essential for testing and demonstrating the AI analytics features.

## What Data is Generated

### 1. Customers (50 records)
- **Customer IDs**: CUST00001 to CUST00050
- **Demographics**: Age (18-78), Gender, Location (10 major US cities)
- **Financial**: Income ranges from $25k to $150k+
- **Behavior**: Registration dates, purchase history, preferred categories
- **Spending Patterns**: Age-based realistic spending amounts

### 2. Purchases (200 records)
- **Products**: 10 different product IDs (PROD001-PROD010)
- **Categories**: Electronics, Fashion, Home & Garden, Books, Sports, Beauty, Automotive
- **Pricing**: Category-based realistic pricing
- **Channels**: Online and in-store purchases
- **Timeline**: Distributed across customer registration periods

### 3. Marketing Campaigns (10 records)
- **Campaign Types**: Email, Social, Display, Search, Influencer
- **Budgets**: Type-based realistic budgets ($500 - $50k)
- **Target Segments**: High Value, Young Adults, Frequent Buyers, etc.
- **Status**: Active, Paused, or Completed based on dates
- **Duration**: 1 week to 2 months

### 4. Campaign Performance (Variable records)
- **Metrics**: Impressions, Clicks, Conversions, Revenue, Cost
- **Calculated KPIs**: CTR, CPC, ROAS automatically calculated
- **Daily Data**: Performance data for each day of campaign duration
- **Realistic Ratios**: Proper conversion funnels

### 5. Customer Segments (5 records)
- **AI Segments**: High Value, At Risk, New, Loyal, Discount Seekers
- **Criteria**: Spending thresholds, activity levels, preferences
- **Sizes**: Realistic segment populations (100-1000 customers)

### 6. Prediction Results (100 records)
- **Types**: Churn prediction, Lifetime Value, Next Purchase timing
- **Probabilities**: Realistic confidence scores (70-100%)
- **Values**: Contextual values based on prediction type
- **Coverage**: Predictions for various customers

## How to Use

### Method 1: Direct Command
```bash
go run seed/main.go
```

### Method 2: Using Makefile
```bash
# Simple seeding (requires MongoDB to be running)
make seed

# Safe seeding (includes MongoDB connectivity check)
make seed-safe
```

### Method 3: Using Script
```bash
chmod +x scripts/seed.sh
./scripts/seed.sh
```

## Prerequisites

1. **MongoDB Running**: Ensure MongoDB is accessible
   ```bash
   # Using Docker
   make docker-run
   
   # Or start local MongoDB
   mongod
   ```

2. **Environment Variables**: Create `.env` file or use defaults
   ```env
   MONGODB_URI=mongodb://localhost:27017
   MONGODB_DBNAME=ai-analytics
   ```

3. **Go Dependencies**: Ensure all dependencies are installed
   ```bash
   go mod tidy
   ```

## Output Example

```
🌱 Starting database seeding...
📊 Generating sample data...
💾 Inserting data into database...
✅ Inserted 50 customers
✅ Inserted 200 purchases
✅ Inserted 10 campaigns
✅ Inserted 157 campaign performances
✅ Inserted 5 customer segments
✅ Inserted 100 prediction results
🎉 Database seeding completed successfully!

📈 Summary:
- Customers: 50
- Purchases: 200
- Campaigns: 10
- Campaign Performances: 157
- Customer Segments: 5
- Prediction Results: 100
```

## Data Relationships

The seeder creates interconnected data:

1. **Customers → Purchases**: Purchases reference valid customer IDs
2. **Campaigns → Performance**: Performance data matches campaign durations
3. **Customers → Predictions**: Predictions reference existing customers
4. **Realistic Timelines**: All dates are logically consistent

## Customization

### Modify Data Volume
Edit the counts in `main()` function:
```go
customers := generateCustomers(100)    // Change from 50 to 100
purchases := generatePurchases(500, customers) // Change from 200 to 500
```

### Add New Categories
Update the category arrays in generation functions:
```go
categories := []string{"Electronics", "Fashion", "YourNewCategory"}
```

### Adjust Spending Patterns
Modify the spending logic in `generateCustomers()`:
```go
if age < 25 {
    totalSpent = float64(100 + rand.Intn(1000)) // Increase spending
}
```

## Database Collections

The seeder populates these MongoDB collections:
- `customers`
- `purchases`
- `campaigns`
- `campaign_performance`
- `customer_segments`
- `predictions`

## Troubleshooting

### MongoDB Connection Issues
```bash
# Check if MongoDB is running
mongosh --eval "db.runCommand('ping')"

# Check connection string in .env
cat .env | grep MONGODB_URI
```

### Permission Issues
```bash
# Make scripts executable
chmod +x scripts/seed.sh
chmod +x scripts/test_auth.sh
```

### Data Already Exists
The seeder will add to existing data. To start fresh:
```bash
# Connect to MongoDB and drop collections
mongosh ai-analytics --eval "
  db.customers.drop();
  db.purchases.drop();
  db.campaigns.drop();
  db.campaign_performance.drop();
  db.customer_segments.drop();
  db.predictions.drop();
"
```

## Integration with AI Analytics

After seeding, you can:

1. **Test Customer Segmentation**: Run K-means clustering on the customer data
2. **Analyze Campaign Performance**: View ROAS, CTR, and other metrics
3. **Predict Customer Behavior**: Test churn and LTV predictions
4. **Dashboard Analytics**: See populated charts and metrics

## Best Practices

1. **Run Once**: Avoid running multiple times unless you want duplicate data
2. **Fresh Start**: Drop collections before re-seeding for clean data
3. **Realistic Data**: The seeder creates realistic patterns for better AI training
4. **Test Environment**: Use seeded data for development and testing only

## Next Steps

After seeding:
1. Start the API server: `make run`
2. Start the frontend: `cd web && npm run dev`
3. Visit http://localhost:3000
4. Login and explore the AI analytics features
5. Generate sample data through the UI for additional testing
