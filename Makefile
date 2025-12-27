.PHONY: setup up down clean build test lint logs status db-migrate db-reset test-api test-db help

# Cluster Management
setup: ## Setup k3d cluster with local registry
	k3d cluster create hikmah-cluster \
		-p "8080:80@loadbalancer" \
		-p "5432:5432@loadbalancer" \
		--agents 1 \
		--registry-create hikmah-registry.localhost:5005

up: ## Launch Tilt (The Dev Loop)
	tilt up

down: ## Stop Tilt and remove resources
	tilt down

clean: ## Delete the entire k3d cluster
	k3d cluster delete hikmah-cluster

# Service Development
build: ## Build all service Docker images
	docker build -t maktba-image -f src/maktba/Dockerfile .
	docker build -t warraq-image -f src/warraq/Dockerfile .
	docker build -t bahith-image -f src/bahith/Dockerfile .
	docker build -t murshid-image -f src/murshid/Dockerfile.

build-maktba: ## Build Maktba service image
	docker build -t maktba-image -f src/maktba/Dockerfile .
	k3d image import maktba-image:latest -c hikmah-cluster

build-warraq: ## Build Warraq service image
	docker build -t warraq-image -f src/warraq/Dockerfile .
	k3d image import warraq-image:latest -c hikmah-cluster

build-bahith: ## Build Bahith service image
	docker build -t bahith-image -f src/bahith/Dockerfile .
	k3d image import bahith-image:latest -c hikmah-cluster

build-murshid: ## Build Murshid service image
	docker build -t murshid-image -f src/murshid/Dockerfile .
	k3d image import murshid-image:latest -c hikmah-cluster

deploy: ## Deploy all services to k3d
	kustomize build deploy/overlays/dev | kubectl apply -f -

test: ## Run all tests
	@echo "Running .NET tests..."
	dotnet test src/maktba/
	@echo "Running Go tests..."
	go test ./src/warraq/...
	go test ./src/bahith/...
	@echo "Running Python tests..."
	python -m pytest src/murshid/

lint: ## Run linting on all services
	@echo "Linting .NET..."
	dotnet format --verify-no-changes src/maktba/
	@echo "Linting Go..."
	golangci-lint run src/warraq/...
	golangci-lint run src/bahith/...
	@echo "Linting Python..."
	flake8 src/murshid/
	mypy src/murshid/

# Database Operations - Updated DB name and password
db-migrate: 
	kubectl exec -it postgres-0 -- psql -U postgres -d hikmah_db -c "\dt"

db-reset: 
	kubectl exec -it postgres-0 -- psql -U postgres -c "DROP DATABASE IF EXISTS hikmah_db; CREATE DATABASE hikmah_db;"

# Testing & Debugging
test-api: ## Test API endpoints
	@echo "Testing Maktba API..."
	curl -f http://localhost:5000/healthz || echo "Maktba health check failed"
	curl -X POST http://localhost:5000/sources -H "Content-Type: application/json" -d '{"title": "Test Book", "type": 0, "description": "Test", "url": "http://test.com"}' || echo "Create source failed"

test-db: ## Test database connectivity
	kubectl exec -it postgres-0 -- psql -U postgres -d hikmah_db -c "SELECT version();" || echo "DB connection failed"

# Utilities
logs: ## Show logs from all services
	kubectl logs -f deployment/maktba &
	kubectl logs -f deployment/warraq &
	kubectl logs -f deployment/bahith &
	kubectl logs -f deployment/murshid &
	kubectl logs -f deployment/gateway &
	kubectl logs -f statefulset/postgres &
	kubectl logs -f statefulset/redpanda

status: ## Show cluster and service status
	@echo "=== Cluster Status ==="
	kubectl get nodes
	@echo ""
	@echo "=== Pods ==="
	kubectl get pods
	@echo ""
	@echo "=== Services ==="
	kubectl get services
	@echo ""
	@echo "=== Recent Events ==="
	kubectl get events --sort-by='.lastTimestamp' | tail -10

port-forward: ## Port forward all services for local access
	kubectl port-forward svc/maktba 5000:80 &
	kubectl port-forward svc/warraq 5001:80 &
	kubectl port-forward svc/bahith 5002:80 &
	kubectl port-forward svc/murshid 5003:80 &
	kubectl port-forward svc/gateway 8080:8080 &
	kubectl port-forward svc/kratos 4433:4433 4434:4434 &
	kubectl port-forward svc/grafana 3000:3000 &
	kubectl port-forward svc/postgres 5432:5432 &
	kubectl port-forward svc/redpanda 9094:9094 8081:8081

shell-maktba: ## Get shell in Maktba pod
	kubectl exec -it deployment/maktba -- bash

shell-db: ## Get shell in PostgreSQL pod
	kubectl exec -it postgres-0 -- psql -U postgres -d hikmah_db

clean-images: ## Clean up Docker images
	docker system prune -f
	docker image prune -f

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
