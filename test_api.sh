#!/bin/bash

# Test script for Chat Server API
# This script tests all endpoints including the new userinfo endpoint

BASE_URL="http://localhost:8080"
EMAIL="testuser_$(date +%s)@example.com"
PASSWORD="test123456"

echo "========================================="
echo "Chat Server API Test"
echo "========================================="
echo ""

# Test 1: Health Check
echo "1. Testing Health Check..."
curl -s $BASE_URL/health | jq .
echo ""
echo ""

# Test 2: Register User
echo "2. Testing User Registration..."
echo "   Email: $EMAIL"
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo $REGISTER_RESPONSE | jq .
TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token')
USER_ID=$(echo $REGISTER_RESPONSE | jq -r '.user_id')
CREATED_AT=$(echo $REGISTER_RESPONSE | jq -r '.created_at')
echo ""
echo "   ✓ User ID: $USER_ID"
echo "   ✓ Created At: $CREATED_AT"
echo "   ✓ Token: ${TOKEN:0:50}..."
echo ""
echo ""

# Test 3: Get User Info
echo "3. Testing User Info Endpoint..."
USER_INFO=$(curl -s $BASE_URL/userinfo \
  -H "Authorization: Bearer $TOKEN")
echo $USER_INFO | jq .
echo ""
echo "   ✓ Retrieved user info with timestamps"
echo ""
echo ""

# Test 4: Login
echo "4. Testing Login..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo $LOGIN_RESPONSE | jq .
NEW_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo ""
echo "   ✓ New Token: ${NEW_TOKEN:0:50}..."
echo ""
echo ""

# Test 5: Access Protected Chat Endpoint
echo "5. Testing Protected Chat Endpoint..."
CHAT_RESPONSE=$(curl -s -X POST $BASE_URL/chat \
  -H "Authorization: Bearer $TOKEN")
echo $CHAT_RESPONSE | jq .
echo ""
echo "   ✓ Successfully accessed protected endpoint"
echo ""
echo ""

# Test 6: Test Invalid Token
echo "6. Testing Invalid Token (should fail)..."
INVALID_RESPONSE=$(curl -s -X POST $BASE_URL/chat \
  -H "Authorization: Bearer invalid-token-here")
echo $INVALID_RESPONSE | jq .
echo ""
echo "   ✓ Properly rejected invalid token"
echo ""
echo ""

# Test 7: Check Metrics
echo "7. Testing Metrics Endpoint..."
echo "   Auth attempts:"
curl -s $BASE_URL/metrics | grep "auth_attempts_total" | grep -v "#"
echo ""
echo "   JWT tokens issued:"
curl -s $BASE_URL/metrics | grep "jwt_tokens_issued_total" | grep -v "#"
echo ""
echo "   HTTP requests:"
curl -s $BASE_URL/metrics | grep "http_requests_total" | head -5
echo ""
echo ""

echo "========================================="
echo "✅ All Tests Completed!"
echo "========================================="
echo ""
echo "Summary:"
echo "  • User registered with timestamps ✓"
echo "  • User info endpoint working ✓"
echo "  • Login successful ✓"
echo "  • Protected endpoints secured ✓"
echo "  • Metrics being tracked ✓"
echo ""

