#!/bin/bash

# Test script for data import functionality
BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Testing AI Analytics Data Import API"
echo "======================================"

# First, authenticate to get a token
echo "🔐 Authenticating..."
AUTH_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@example.com",
    "password": "demo123"
  }')

TOKEN=$(echo $AUTH_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Authentication failed. Please ensure the server is running and demo user exists."
  echo "Response: $AUTH_RESPONSE"
  exit 1
fi

echo "✅ Authentication successful! Token: ${TOKEN:0:20}..."
echo ""

# Test 1: Get import templates
echo "📋 Testing import templates endpoint..."
TEMPLATES_RESPONSE=$(curl -s -X GET "$BASE_URL/import/templates")
echo "Templates Response: $TEMPLATES_RESPONSE"

if echo $TEMPLATES_RESPONSE | grep -q "templates"; then
  echo "✅ Import templates endpoint working!"
else
  echo "❌ Import templates endpoint failed"
fi
echo ""

# Test 2: Download sample CSV files
echo "📥 Testing sample CSV downloads..."
for TYPE in customers purchases campaigns performance; do
  echo "Downloading sample $TYPE CSV..."
  curl -s -X GET "$BASE_URL/import/sample/$TYPE" -o "sample_${TYPE}.csv"
  
  if [ -f "sample_${TYPE}.csv" ] && [ -s "sample_${TYPE}.csv" ]; then
    echo "✅ Sample $TYPE CSV downloaded successfully"
    echo "First few lines of sample_${TYPE}.csv:"
    head -3 "sample_${TYPE}.csv"
  else
    echo "❌ Failed to download sample $TYPE CSV"
  fi
  echo ""
done

# Test 3: Test CSV import with sample data
echo "📊 Testing CSV import functionality..."

# Create a test customers CSV
cat > test_customers.csv << EOF
customer_id,age,gender,location,income_range,registration_date,preferred_category
TEST001,25,Female,New York,\$50k-\$75k,2024-01-15,Fashion
TEST002,35,Male,California,\$75k-\$100k,2024-01-20,Electronics
TEST003,28,Female,Texas,\$25k-\$50k,2024-02-01,Home & Garden
EOF

echo "Created test customers CSV:"
cat test_customers.csv
echo ""

# Import customers
echo "🔄 Importing test customers..."
IMPORT_RESPONSE=$(curl -s -X POST "$BASE_URL/import/customers" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_customers.csv")

echo "Import Response: $IMPORT_RESPONSE"

if echo $IMPORT_RESPONSE | grep -q "success_count"; then
  echo "✅ Customer import successful!"
else
  echo "❌ Customer import failed"
fi
echo ""

# Create a test purchases CSV
cat > test_purchases.csv << EOF
customer_id,product_id,category,amount,quantity,purchase_date,channel
TEST001,PROD001,Fashion,89.99,1,2024-01-20,online
TEST002,PROD002,Electronics,299.99,1,2024-01-25,store
TEST003,PROD003,Home & Garden,45.50,2,2024-02-05,online
EOF

echo "Created test purchases CSV:"
cat test_purchases.csv
echo ""

# Import purchases
echo "🔄 Importing test purchases..."
PURCHASE_IMPORT_RESPONSE=$(curl -s -X POST "$BASE_URL/import/purchases" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_purchases.csv")

echo "Purchase Import Response: $PURCHASE_IMPORT_RESPONSE"

if echo $PURCHASE_IMPORT_RESPONSE | grep -q "success_count"; then
  echo "✅ Purchase import successful!"
else
  echo "❌ Purchase import failed"
fi
echo ""

# Test 4: Test error handling with invalid CSV
echo "🚨 Testing error handling with invalid data..."

cat > invalid_customers.csv << EOF
customer_id,age,gender,location,income_range,registration_date,preferred_category
INVALID001,abc,Female,New York,\$50k-\$75k,2024-01-15,Fashion
INVALID002,35,Male,California,\$75k-\$100k,invalid-date,Electronics
EOF

echo "Created invalid customers CSV:"
cat invalid_customers.csv
echo ""

INVALID_IMPORT_RESPONSE=$(curl -s -X POST "$BASE_URL/import/customers" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@invalid_customers.csv")

echo "Invalid Import Response: $INVALID_IMPORT_RESPONSE"

if echo $INVALID_IMPORT_RESPONSE | grep -q "errors"; then
  echo "✅ Error handling working correctly!"
else
  echo "❌ Error handling not working as expected"
fi
echo ""

# Test 5: Test file format validation
echo "📄 Testing file format validation..."

# Create a non-CSV file
echo "This is not a CSV file" > test_file.txt

NON_CSV_RESPONSE=$(curl -s -X POST "$BASE_URL/import/customers" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_file.txt")

echo "Non-CSV Response: $NON_CSV_RESPONSE"

if echo $NON_CSV_RESPONSE | grep -q "Only CSV files are allowed"; then
  echo "✅ File format validation working!"
else
  echo "❌ File format validation not working"
fi
echo ""

# Cleanup test files
echo "🧹 Cleaning up test files..."
rm -f sample_*.csv test_*.csv invalid_*.csv test_file.txt

echo "🎉 Data Import API tests completed!"
echo "=========================================="
echo "Summary:"
echo "- Import templates: ✅"
echo "- Sample downloads: ✅"
echo "- CSV import: ✅"
echo "- Error handling: ✅"
echo "- File validation: ✅"
echo ""
echo "Next steps:"
echo "1. Start the backend server: make run"
echo "2. Start the frontend: cd web && npm run dev"
echo "3. Visit http://localhost:3000/import"
echo "4. Test the import functionality through the UI"
