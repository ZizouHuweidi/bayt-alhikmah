#!/bin/bash

# Configuration
GATEWAY_URL="http://localhost:8080"
echo "Waiting for services to be up..."
sleep 5

# 1. Test Gateway Health (or just root)
echo "Testing Gateway connectivity..."
curl -I "$GATEWAY_URL/v1/sources" || echo "Gateway might be starting..."

# 2. Test Kratos (Auth) via Gateway
echo "Testing Auth Endpoint (Public Config)..."
curl -s "$GATEWAY_URL/auth/ui/welcome" | grep "welcome" && echo "Auth UI Reachable" || echo "Auth Failed"

# 3. Create a Source (Maktba) via Gateway
echo "Creating a Source..."
RESPONSE=$(curl -s -X POST "$GATEWAY_URL/v1/sources" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The House of Wisdom",
    "type": 0,
    "description": "A book about the history of Bayt al-Hikmah",
    "url": "https://example.com/book"
  }')

echo "Response: $RESPONSE"

# Extract ID (simple grep or jq if avail)
ID=$(echo $RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

if [ -z "$ID" ]; then
    echo "Failed to create source."
else
    echo "Source Created with ID: $ID"
    # 4. Get Source
    echo "Retrieving Source..."
    curl -s "$GATEWAY_URL/v1/sources/$ID"
fi
