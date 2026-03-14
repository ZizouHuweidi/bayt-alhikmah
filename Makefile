.PHONY: help up down restart logs clean build test lint fmt migrate db deps frontend status health dev

# Default target
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "$(BLUE)%-25s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ============================================================================
# Docker Compose Operations
# ============================================================================

up:  ## Start all core services (postgres, redis, redpanda, auth, maktaba)
	@echo "$(GREEN)Starting Bayt al Hikmah core services...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Waiting for services to be healthy...$(NC)"
	@sleep 10
	@make health

down:  ## Stop all services
	@echo "$(YELLOW)Stopping all services...$(NC)"
	docker-compose down

restart:  ## Restart all services
	@echo "$(YELLOW)Restarting all services...$(NC)"
	docker-compose restart

logs:  ## View logs from all services
	docker-compose logs -f

logs-%:  ## View logs from specific service (e.g., make logs-maktaba)
	docker-compose logs -f $*

ps:  ## Show running containers
	docker-compose ps

status:  ## Show status of all services
	@echo "$(BLUE)=== Service Status ===$(NC)"
	@docker-compose ps
	@echo ""
	@make health

health:  ## Check health of services
	@echo "$(BLUE)=== Health Checks ===$(NC)"
	@curl -s http://localhost:8080/health 2>/dev/null && echo "$(GREEN)✓ Maktaba: OK$(NC)" || echo "$(RED)✗ Maktaba: DOWN$(NC)"
	@curl -s http://localhost:4433/health/ready 2>/dev/null && echo "$(GREEN)✓ Kratos: OK$(NC)" || echo "$(RED)✗ Kratos: DOWN$(NC)"
	@docker-compose exec -T postgres pg_isready -U postgres 2>/dev/null && echo "$(GREEN)✓ PostgreSQL: OK$(NC)" || echo "$(RED)✗ PostgreSQL: DOWN$(NC)"
	@docker-compose exec -T redis redis-cli ping 2>/dev/null | grep -q PONG && echo "$(GREEN)✓ Redis: OK$(NC)" || echo "$(RED)✗ Redis: DOWN$(NC)"

clean:  ## Remove all containers, volumes, and networks
	@echo "$(RED)WARNING: This will delete all data!$(NC)"
	@read -p "Are you sure? [y/N] " confirm && [ $$confirm = y ] && docker-compose down -v || echo "Cancelled"

clean-volumes:  ## Remove all volumes (WARNING: deletes all data)
	@echo "$(RED)WARNING: This will delete all data permanently!$(NC)"
	@read -p "Are you sure? [y/N] " confirm && [ $$confirm = y ] && docker-compose down -v && docker volume rm $$(docker volume ls -q -f name=bayt-alhikmah) 2>/dev/null || echo "Cancelled"

rebuild: clean-volumes build-maktaba up migrate  ## Clean rebuild with fresh data

# ============================================================================
# Build Commands
# ============================================================================

build: build-maktaba  ## Build all services

build-maktaba:  ## Build maktaba service image
	@echo "$(BLUE)Building maktaba service...$(NC)"
	cd src/maktaba && docker build -t bayt-alhikmah/maktaba:latest .

build-bahith:  ## Build bahith service image
	@echo "$(BLUE)Building bahith service...$(NC)"
	cd src/bahith && docker build -t bayt-alhikmah/bahith:latest .

build-murshid:  ## Build murshid service image
	@echo "$(BLUE)Building murshid service...$(NC)"
	cd src/murshid && docker build -t bayt-alhikmah/murshid:latest .

build-all: build-maktaba build-bahith build-murshid  ## Build all service images

# ============================================================================
# Testing
# ============================================================================

test: test-maktaba  ## Run all tests

test-maktaba:  ## Run maktaba tests
	@echo "$(BLUE)Running maktaba tests...$(NC)"
	cd src/maktaba && go test ./...

test-bahith:  ## Run bahith tests
	@echo "$(BLUE)Running bahith tests...$(NC)"
	cd src/bahith && pytest

test-murshid:  ## Run murshid tests
	@echo "$(BLUE)Running murshid tests...$(NC)"
	cd src/murshid && pytest

test-all: test-maktaba test-bahith test-murshid  ## Run all tests

# ============================================================================
# Linting & Formatting
# ============================================================================

lint: lint-maktaba  ## Lint all code

lint-maktaba:  ## Lint maktaba code
	@echo "$(BLUE)Linting maktaba...$(NC)"
	cd src/maktaba && golangci-lint run

lint-bahith:  ## Lint bahith code
	@echo "$(BLUE)Linting bahith...$(NC)"
	cd src/bahith && ruff check . && black --check .

lint-murshid:  ## Lint murshid code
	@echo "$(BLUE)Linting murshid...$(NC)"
	cd src/murshid && ruff check . && black --check .

lint-all: lint-maktaba lint-bahith lint-murshid  ## Lint all code

fmt: fmt-maktaba  ## Format all code

fmt-maktaba:  ## Format maktaba code
	@echo "$(BLUE)Formatting maktaba...$(NC)"
	cd src/maktaba && gofmt -s -w .

fmt-bahith:  ## Format bahith code
	@echo "$(BLUE)Formatting bahith...$(NC)"
	cd src/bahith && ruff check --fix . && black .

fmt-murshid:  ## Format murshid code
	@echo "$(BLUE)Formatting murshid...$(NC)"
	cd src/murshid && ruff check --fix . && black .

fmt-all: fmt-maktaba fmt-bahith fmt-murshid  ## Format all code

# ============================================================================
# Database & Migrations
# ============================================================================


migrate:  ## Run maktaba database migrations inside container
	@echo "$(BLUE)Running migrations...$(NC)"
	@docker-compose exec -T postgres psql -U postgres -d maktaba -f /migrations/000001_initial.up.sql 2>/dev/null || \
	(docker-compose exec postgres psql -U postgres -c "CREATE DATABASE maktaba;" && \
	 docker-compose exec -T postgres psql -U postgres -d maktaba -f /migrations/000001_initial.up.sql)

migrate-down:  ## Rollback maktaba migrations
	@echo "$(YELLOW)Rolling back maktaba migrations...$(NC)"
	@docker-compose exec -T postgres psql -U postgres -d maktaba -f /migrations/000001_initial.down.sql 2>/dev/null || echo "$(RED)Rollback failed or no migrations to rollback$(NC)"

db-create:  ## Create maktaba database
	@echo "$(BLUE)Creating maktaba database...$(NC)"
	@docker-compose exec postgres psql -U postgres -c "CREATE DATABASE maktaba;" 2>/dev/null && echo "$(GREEN)Database created$(NC)" || echo "$(YELLOW)Database already exists$(NC)"

db-drop:  ## Drop maktaba database
	@echo "$(RED)Dropping maktaba database...$(NC)"
	@docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS maktaba;" && echo "$(GREEN)Database dropped$(NC)"

db-reset: db-drop db-create migrate  ## Reset maktaba database (drop, create, migrate)

db-shell:  ## Open PostgreSQL shell
	docker-compose exec postgres psql -U postgres

db-shell-maktaba:  ## Open PostgreSQL shell for maktaba database
	docker-compose exec postgres psql -U postgres -d maktaba

sqlc-gen:  ## Generate Go code from SQL queries
	@echo "$(BLUE)Generating SQLC code...$(NC)"
	cd src/maktaba && sqlc generate

# ============================================================================
# Dependencies
# ============================================================================

deps: deps-maktaba deps-frontend  ## Install all dependencies

deps-maktaba:  ## Install maktaba dependencies
	@echo "$(BLUE)Installing maktaba dependencies...$(NC)"
	cd src/maktaba && go mod download && go mod tidy

deps-bahith:  ## Install bahith dependencies
	@echo "$(BLUE)Installing bahith dependencies...$(NC)"
	cd src/bahith && uv sync

deps-murshid:  ## Install murshid dependencies
	@echo "$(BLUE)Installing murshid dependencies...$(NC)"
	cd src/murshid && uv sync

deps-frontend:  ## Install frontend dependencies
	@echo "$(BLUE)Installing frontend dependencies...$(NC)"
	cd frontend && bun install

deps-all: deps-maktaba deps-bahith deps-murshid deps-frontend  ## Install all dependencies

# ============================================================================
# Frontend
# ============================================================================

frontend-dev:  ## Start frontend development server
	@echo "$(GREEN)Starting frontend development server...$(NC)"
	cd frontend && bun run dev

frontend-build:  ## Build frontend for production
	@echo "$(BLUE)Building frontend...$(NC)"
	cd frontend && bun run build

frontend-lint:  ## Lint frontend code
	@echo "$(BLUE)Linting frontend...$(NC)"
	cd frontend && bun run lint

# ============================================================================
# Services Management
# ============================================================================

redis-cli:  ## Open Redis CLI
	docker-compose exec redis redis-cli

kafka-topics:  ## List Kafka topics
	docker-compose exec redpanda rpk topic list

kratos-ready:  ## Check if Kratos is ready
	@curl -s http://localhost:4433/health/ready 2>/dev/null && echo "$(GREEN)Kratos is ready$(NC)" || echo "$(RED)Kratos is not ready$(NC)"

# ============================================================================
# Quick URLs
# ============================================================================

open-landing:  ## Open landing page
	@xdg-open http://localhost:3000 2>/dev/null || open http://localhost:3000 || echo "Open http://localhost:3000"

open-login:  ## Open login page
	@xdg-open http://localhost:3000/login 2>/dev/null || open http://localhost:3000/login || echo "Open http://localhost:3000/login"

open-dashboard:  ## Open dashboard
	@xdg-open http://localhost:3000/dashboard 2>/dev/null || open http://localhost:3000/dashboard || echo "Open http://localhost:3000/dashboard"

open-api:  ## Open Maktaba API
	@xdg-open http://localhost:8080 2>/dev/null || open http://localhost:8080 || echo "Open http://localhost:8080"

open-pgadmin:  ## Open pgAdmin
	@xdg-open http://localhost:5050 2>/dev/null || open http://localhost:5050 || echo "Open http://localhost:5050"

# ============================================================================
# Development Workflow
# ============================================================================

dev: deps up  ## Full development setup (install deps, start services, run migrations)
	@echo "$(GREEN)Waiting for services to start...$(NC)"
	@sleep 5
	@make migrate
	@echo ""
	@echo "$(GREEN)==============================================$(NC)"
	@echo "$(GREEN)Development environment is ready!$(NC)"
	@echo "$(GREEN)==============================================$(NC)"
	@echo ""
	@echo "$(BLUE)Services:$(NC)"
	@echo "  Frontend:   http://localhost:3000"
	@echo "  Maktaba:    http://localhost:8080"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  pgAdmin:    http://localhost:5050"
	@echo ""
	@echo "$(BLUE)Next steps:$(NC)"
	@echo "  1. Run 'make frontend-dev' to start the frontend"
	@echo "  2. Visit http://localhost:3000 to see the landing page"
	@echo "  3. Register at http://localhost:3000/registration"
	@echo ""

dev-full: dev frontend-dev  ## Full dev setup + start frontend

# ============================================================================
# Auth Testing
# ============================================================================

test-auth-flow:  ## Test authentication flow
	@echo "$(BLUE)Testing auth setup...$(NC)"
	@echo ""
	@echo "1. Checking Kratos health..."
	@curl -s http://localhost:4433/health/ready 2>/dev/null | grep -q "ok" && echo "$(GREEN)   ✓ Kratos is healthy$(NC)" || echo "$(RED)   ✗ Kratos is not healthy$(NC)"
	@echo ""
	@echo "2. Checking Maktaba health..."
	@curl -s http://localhost:8080/health 2>/dev/null | grep -q "healthy" && echo "$(GREEN)   ✓ Maktaba is healthy$(NC)" || echo "$(RED)   ✗ Maktaba is not healthy$(NC)"
	@echo ""
	@echo "3. Creating test flow..."
	@curl -s -X GET http://localhost:4433/self-service/registration/api 2>/dev/null | grep -q "id" && echo "$(GREEN)   ✓ Can create registration flow$(NC)" || echo "$(RED)   ✗ Cannot create registration flow$(NC)"
	@echo ""
	@echo "$(GREEN)Auth system check complete!$(NC)"
	@echo "Visit http://localhost:3000 to test the full flow."
