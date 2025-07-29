#!/bin/bash

# Test Authentication API
BASE_URL="http://localhost:8080"

echo "=== Testing Authentication System ==="
echo

# Test user registration
echo "1. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }')

echo "Register Response: $REGISTER_RESPONSE"
echo

# Extract token from registration response
TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "Registration failed, trying login instead..."
  
  # Test user login
  echo "2. Testing user login..."
  LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
      "email": "test@example.com",
      "password": "password123"
    }')
  
  echo "Login Response: $LOGIN_RESPONSE"
  echo
  
  # Extract token from login response
  TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
fi

if [ -n "$TOKEN" ]; then
  echo "Token obtained: ${TOKEN:0:50}..."
  echo
  
  # Test protected route - get current user
  echo "3. Testing protected route - /api/auth/me..."
  ME_RESPONSE=$(curl -s -X GET "$BASE_URL/api/auth/me" \
    -H "Authorization: Bearer $TOKEN")
  
  echo "Me Response: $ME_RESPONSE"
  echo
  
  # Test protected route - profile
  echo "4. Testing protected route - /api/protected/profile..."
  PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/protected/profile" \
    -H "Authorization: Bearer $TOKEN")
  
  echo "Profile Response: $PROFILE_RESPONSE"
  echo
  
  # Test protected route - dashboard
  echo "5. Testing protected route - /api/protected/dashboard..."
  DASHBOARD_RESPONSE=$(curl -s -X GET "$BASE_URL/api/protected/dashboard" \
    -H "Authorization: Bearer $TOKEN")
  
  echo "Dashboard Response: $DASHBOARD_RESPONSE"
  echo
  
  # Test unauthorized access
  echo "6. Testing unauthorized access..."
  UNAUTHORIZED_RESPONSE=$(curl -s -X GET "$BASE_URL/api/protected/profile")
  
  echo "Unauthorized Response: $UNAUTHORIZED_RESPONSE"
  echo
else
  echo "Failed to obtain token. Check if the server is running and MongoDB is connected."
fi

echo "=== Authentication tests completed ==="
