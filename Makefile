# Load environment variables from .env file
include .env
export

# ====================================================================================
# VARIABLES
# ====================================================================================

# This is the connection string for the migrate tool running inside Docker.
# It uses the DB_HOST service name ('db') to connect to the postgres container.
POSTGRES_URL := postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

# Absolute path to the migrations directory
MIGRATIONS_ROOT := $(shell pwd)/migrations

# Docker network name from docker-compose.yml
NETWORK := bayt-alhikmah_bayt-alhikmah-net

# ====================================================================================
# HELP
# ====================================================================================

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_%\/-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ====================================================================================
# DEVELOPMENT WORKFLOW
# ====================================================================================

.PHONY: run
run: ## Run the application locally (without live reload)
	@go run cmd/api/main.go

.PHONY: watch
watch: ## Run the application with Air for live reloading
	@if ! command -v air &> /dev/null; then \
		echo "Air not found. Installing..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@air

.PHONY: docker/up
docker/up: ## Start all services with Docker Compose
	@docker compose up --build -d

.PHONY: docker/down
docker/down: ## Stop all services with Docker Compose
	@docker compose down

.PHONY: docker/logs
docker/logs: ## Follow logs for the app service
	@docker compose logs -f app

# ====================================================================================
# DATABASE & MIGRATIONS
# ====================================================================================

.PHONY: db/generate
db/generate: ## Generate Go code from SQL queries using SQLc
	@if ! command -v sqlc &> /dev/null; then \
		echo "sqlc not found. Installing..."; \
		go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
	fi
	@sqlc generate

.PHONY: migrate/new
migrate/new: ## Create new up/down migration files. Usage: make migrate/new name=create_books_table
	@docker run --rm -v $(MIGRATIONS_ROOT):/migrations migrate/migrate create -ext sql -dir /migrations -seq $(name)

.PHONY: migrate/up
migrate/up: ## Apply all available 'up' migrations
	@docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database '$(POSTGRES_URL)' up

.PHONY: migrate/down
migrate/down: ## Revert the last migration
	@docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database '$(POSTGRES_URL)' down 1

.PHONY: migrate/down/all
migrate/down/all: ## Revert all migrations
	@docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database '$(POSTGRES_URL)' down -all

.PHONY: db/connect
db/connect: ## Connect to the running PostgreSQL database using psql
	@docker compose exec -u postgres db psql -d ${DB_NAME}

# ====================================================================================
# QUALITY & TESTING
# ====================================================================================

.PHONY: build
build: ## Build the application binary
	@echo "Building..."
	@go build -o build/app cmd/api/main.go

.PHONY: test
test: ## Run all Go tests
	@echo "Running tests..."
	@go test ./... -v

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum files
	@go mod tidy
	@go mod verify

.PHONY: clean
clean: ## Clean the build binary
	@echo "Cleaning..."
	@rm -f build/app
