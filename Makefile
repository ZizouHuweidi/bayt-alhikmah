.PHONY: setup up down clean build test lint logs status db-migrate db-reset test-api test-db help port-forward

# Cluster Management
# Images to pull
IMAGES = \
	grafana/grafana:11.4.0 \
	grafana/tempo:2.6.1 \
	prom/prometheus:v3.0.1 \
	grafana/loki:3.3.2 \
	oryd/kratos:v1.3.1 \
	postgres:18-alpine \
	docker.redpanda.com/redpandadata/redpanda:v24.2.5 \
	devopsfaith/krakend:2.7.0 \
	mcr.microsoft.com/dotnet/aspnet:10.0 \
	mcr.microsoft.com/dotnet/sdk:10.0

pull-images: ## Pull all required images sequentially
	@echo "Pulling images one by one..."
	@for img in $(IMAGES); do \
		echo "Pulling $$img..."; \
		docker pull $$img; \
	done

setup: ## Setup k3d cluster with local registry and pre-pulled images
	$(MAKE) pull-images
	k3d cluster create hikmah-cluster \
		--agents 1 \
		--registry-create hikmah-registry.localhost:5005
	@echo "Importing images into cluster (this may take a minute)..."
	k3d image import $(IMAGES) -c hikmah-cluster --verbose

up: ## Launch Tilt (The Dev Loop)
	tilt up

down: ## Stop Tilt and remove resources
	tilt down

clean: ## Delete the entire k3d cluster
	k3d cluster delete hikmah-cluster

# Service Development
build: build-maktba build-gateway build-warraq build-bahith build-murshid ## Build all service images

build-maktba: ## Build Maktba service image
	docker build -t maktba-image -f src/maktba/Dockerfile .
	k3d image import maktba-image:latest -c hikmah-cluster

build-gateway: ## Build Gateway image
	docker build -t gateway-image -f src/madkhal/Dockerfile src/madkhal
	k3d image import gateway-image:latest -c hikmah-cluster

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

lint: ## Run linting on all services
	@echo "Linting .NET..."
	dotnet format --verify-no-changes src/maktba/
	@echo "Linting Go..."
	golangci-lint run src/warraq/...

# Database Operations
db-migrate: ## Show tables in database
	kubectl exec -it postgres-0 -- psql -U postgres -d hikmah_db -c "\dt"

db-reset: ## Reset the hikmah_db database
	kubectl exec -it postgres-0 -- psql -U postgres -c "DROP DATABASE IF EXISTS hikmah_db; CREATE DATABASE hikmah_db;"

db-shell: ## Get shell in PostgreSQL pod
	kubectl exec -it postgres-0 -- psql -U postgres -d hikmah_db

# Testing & Debugging
test-api: ## Run verification script
	./scripts/verify_api.sh

test-db: ## Test database connectivity
	kubectl exec -it postgres-0 -- psql -U postgres -d hikmah_db -c "SELECT version();" || echo "DB connection failed"

# Utilities
logs: ## Show logs from all core services
	kubectl logs -l app=gateway -f &
	kubectl logs -l app=maktba -f &
	kubectl logs -l app=kratos -f

status: ## Show cluster and service status
	@echo "=== Cluster Status ==="
	kubectl get nodes
	@echo ""
	@echo "=== Pods ==="
	kubectl get pods
	@echo ""
	@echo "=== Services ==="
	kubectl get services

port-forward: ## Port forward all services for local access
	kubectl port-forward svc/gateway 8080:8080 &
	kubectl port-forward svc/maktba 5000:80 &
	kubectl port-forward svc/kratos 4433:4433 4434:4434 &
	kubectl port-forward svc/grafana 3000:3000 &
	kubectl port-forward svc/prometheus-server 9090:9090 &
	kubectl port-forward svc/loki 3100:3100 &
	kubectl port-forward svc/tempo 3200:3200 4317:4317 &
	kubectl port-forward svc/postgres 5432:5432 &
	kubectl port-forward svc/redpanda 9092:9092 9644:9644

clean-images: ## Clean up Docker images
	docker system prune -f

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
