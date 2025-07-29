#!/bin/bash

# Test script for authentication endpoints
BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Testing AI Analytics Authentication API"
echo "=========================================="

# Test user registration
echo "📝 Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "test123"
  }')

echo "Registration Response: $REGISTER_RESPONSE"

# Extract token from registration response
TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
  echo "✅ Registration successful! Token: ${TOKEN:0:20}..."
else
  echo "❌ Registration failed"
  exit 1
fi

echo ""

# Test user login
echo "🔐 Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "test123"
  }')

echo "Login Response: $LOGIN_RESPONSE"

# Extract token from login response
LOGIN_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$LOGIN_TOKEN" ]; then
  echo "✅ Login successful! Token: ${LOGIN_TOKEN:0:20}..."
else
  echo "❌ Login failed"
  exit 1
fi

echo ""

# Test protected endpoint
echo "🔒 Testing protected endpoint..."
ME_RESPONSE=$(curl -s -X GET "$BASE_URL/auth/me" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "Me Response: $ME_RESPONSE"

if echo $ME_RESPONSE | grep -q "user"; then
  echo "✅ Protected endpoint access successful!"
else
  echo "❌ Protected endpoint access failed"
  exit 1
fi

echo ""

# Test analytics dashboard
echo "📊 Testing analytics dashboard..."
DASHBOARD_RESPONSE=$(curl -s -X GET "$BASE_URL/analytics/dashboard" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "Dashboard Response: $DASHBOARD_RESPONSE"

if echo $DASHBOARD_RESPONSE | grep -q "dashboard"; then
  echo "✅ Analytics dashboard access successful!"
else
  echo "❌ Analytics dashboard access failed"
fi

echo ""

# Test sample data generation
echo "🎲 Testing sample data generation..."
SAMPLE_RESPONSE=$(curl -s -X POST "$BASE_URL/analytics/sample-data")

echo "Sample Data Response: $SAMPLE_RESPONSE"

if echo $SAMPLE_RESPONSE | grep -q "sample_data"; then
  echo "✅ Sample data generation successful!"
else
  echo "❌ Sample data generation failed"
fi

echo ""
echo "🎉 Authentication tests completed!"
echo "=========================================="
echo "Next steps:"
echo "1. Start the backend server: make run"
echo "2. Start the frontend: cd web && npm run dev"
echo "3. Visit http://localhost:3000/login"
echo "4. Use credentials: test@example.com / test123"
