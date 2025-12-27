#!/bin/bash

# API Testing Script for Bayt al-Hikmah
# Tests the core API endpoints via the Gateway

set -e

# Configuration
GATEWAY_URL="http://localhost:8080"

echo "🧪 Testing Bayt al-Hikmah APIs via Gateway ($GATEWAY_URL)..."

# Function to test endpoint
test_endpoint() {
    local url=$1
    local method=${2:-"GET"}
    local data=${3:-""}
    local description=$4
    
    echo "🔍 $description"
    
    local cmd="curl -s"
    if [ "$method" = "POST" ]; then
        cmd="$cmd -X POST -H \"Content-Type: application/json\""
        if [ -n "$data" ]; then
            cmd="$cmd -d '$data'"
        fi
    fi
    
    response=$(eval "$cmd \"$url\"" || echo "ERROR")
    
    if [ "$response" = "ERROR" ] || [ -z "$response" ]; then
        echo "❌ Failed"
        return 1
    else
        echo "✅ Success"
        echo "$response" | head -c 200
        echo -e "\n---"
        return 0
    fi
}

# 1. Test Gateway Health via Kratos Ready Check
test_endpoint "$GATEWAY_URL/auth/health/ready" "GET" "" "Gateway -> Kratos Health"

# 2. Maktba API via Gateway
echo "📡 Testing Maktba API via Gateway..."

SOURCE_DATA='{
    "title": "The House of Wisdom",
    "type": 0,
    "description": "Verified at $(date)",
    "url": "http://example.com/verify"
}'

response=$(curl -s -X POST "$GATEWAY_URL/v1/sources" \
    -H "Content-Type: application/json" \
    -d "$SOURCE_DATA")

if [ -n "$response" ] && [[ "$response" == *"id"* ]]; then
    SOURCE_ID=$(echo "$response" | grep -oP '"id":"\K[^"]+')
    echo "✅ Created source with ID: $SOURCE_ID"
    
    # Test retrieving
    test_endpoint "$GATEWAY_URL/v1/sources/$SOURCE_ID" "GET" "" "Get source by ID"
else
    echo "❌ Failed to create source via Gateway. Response: $response"
fi

# 3. Kratos Auth Flow via Gateway
echo ""
echo "🔐 Testing Kratos Auth Flow via Gateway..."

# Get registration flow
FLOW_JSON=$(curl -s "$GATEWAY_URL/auth/self-service/registration/api")
FLOW_ID=$(echo "$FLOW_JSON" | grep -oP '"id":"\K[^"]+')

if [ -n "$FLOW_ID" ]; then
    echo "✅ Successfully initiated registration flow: $FLOW_ID"
else
    echo "❌ Failed to initiate auth flow via Gateway"
fi

echo ""
echo "✅ API testing complete!"
echo ""
echo "💡 If tests failed, make sure 'tilt up' is running or use 'make port-forward'."
