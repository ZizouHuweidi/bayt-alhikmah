#!/bin/bash

# API Testing Script for Bayt al-Hikmah
# Tests the core API endpoints

set -e

# Configuration
GATEWAY_URL="http://localhost:8080"
MAKTBA_URL="http://localhost:5000"

echo "🧪 Testing Bayt al-Hikmah APIs..."

# Function to test endpoint
test_endpoint() {
    local url=$1
    local method=${2:-"GET"}
    local data=${3:-""}
    local description=$4
    
    echo "🔍 $description"
    
    if [ "$method" = "POST" ] && [ -n "$data" ]; then
        response=$(curl -s -X POST "$url" -H "Content-Type: application/json" -d "$data" || echo "ERROR")
    else
        response=$(curl -s -f "$url" || echo "ERROR")
    fi
    
    if [ "$response" = "ERROR" ]; then
        echo "❌ Failed"
        return 1
    else
        echo "✅ Success"
        echo "$response" | head -3
        echo "---"
        return 0
    fi
}

# Test direct Maktba service first
echo "📡 Testing Maktba service directly..."
test_endpoint "$MAKTBA_URL/healthz" "GET" "" "Health check"

# Create a test source via Maktba
SOURCE_DATA='{
    "title": "The House of Wisdom",
    "type": 0,
    "description": "How Arabic Science Saved Ancient Knowledge",
    "url": "http://example.com/book"
}'

response=$(curl -s -X POST "$MAKTBA_URL/sources" \
    -H "Content-Type: application/json" \
    -d "$SOURCE_DATA" || echo "ERROR")

if [ "$response" != "ERROR" ]; then
    SOURCE_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 || echo "")
    echo "✅ Created source with ID: $SOURCE_ID"
    
    # Test retrieving the source
    if [ -n "$SOURCE_ID" ]; then
        test_endpoint "$MAKTBA_URL/sources/$SOURCE_ID" "GET" "" "Get source by ID"
    fi
else
    echo "❌ Failed to create source via Maktba"
fi


# Test Gateway connectivity
echo ""
echo "🌐 Testing Gateway connectivity..."
test_endpoint "$GATEWAY_URL/v1/sources" "GET" "" "Gateway sources endpoint"

if [ -n "$SOURCE_ID" ]; then
    test_endpoint "$GATEWAY_URL/v1/sources/$SOURCE_ID" "GET" "" "Gateway get source"
fi

# Test Kratos auth endpoints
echo ""
echo "🔐 Testing Auth endpoints..."
test_endpoint "http://localhost:4434/.well-known/ory/kratos/public/" "GET" "" "Kratos public config"
test_endpoint "http://localhost:4433/self-service/registration/browser" "GET" "" "Kratos registration"

echo ""
echo "✅ API testing complete!"
echo ""
echo "💡 If tests failed, check:"
echo "   - kubectl get pods (are services running?)"
echo "   - make logs (check service logs)"
echo "   - make port-forward (are ports forwarded?)"
