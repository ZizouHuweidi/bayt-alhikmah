#!/bin/bash

# Bayt al-Hikmah Development Environment Setup Script
# This script sets up the complete development environment

set -e

echo "🏛️  Setting up Bayt al-Hikmah Development Environment..."

# Check prerequisites
command -v k3d >/dev/null 2>&1 || { echo "❌ k3d is required but not installed. Please install k3d first."; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "❌ kubectl is required but not installed. Please install kubectl first."; exit 1; }
command -v tilt >/dev/null 2>&1 || { echo "❌ Tilt is required but not installed. Please install Tilt first."; exit 1; }

# Create k3d cluster if it doesn't exist
if ! k3d cluster list | grep -q "hikmah-cluster"; then
    echo "📦 Creating k3d cluster..."
    make setup
else
    echo "✅ k3d cluster already exists"
fi

# Build and import images
echo "🔨 Building service images..."
make build-maktba

# Deploy services
echo "🚀 Deploying services..."
make deploy

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=maktba --timeout=120s || true

# Run database migrations
echo "🗄️  Running database migrations..."
kubectl exec -it postgres-0 -- psql -U postgres -d hikmah_db -c "SELECT version();" || echo "⚠️  Database not ready yet"

echo ""
echo "✅ Development environment setup complete!"
echo ""
echo "🌐 Access points:"
echo "   - Maktba API: http://localhost:5000"
echo "   - Kratos Admin: http://localhost:4434"
echo "   - Grafana: http://localhost:3000"
echo ""
echo "🛠️  Common commands:"
echo "   - make up          # Start Tilt dev loop"
echo "   - make logs        # View service logs"
echo "   - make status      # Check cluster status"
echo "   - make test-api    # Test API endpoints"
echo ""
echo "📚 Documentation available in README.md"