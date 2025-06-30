# Makefile for bayt-alhikmah

.DEFAULT_GOAL := help

# ==============================================================================
# Help
# ==============================================================================

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ==============================================================================
# Docker Compose Commands
# ==============================================================================

.PHONY: up
up: ## Start all services in the background using Docker Compose
	@echo "Starting services..."
	@docker compose up -d --build

.PHONY: down
down: ## Stop all services started with Docker Compose
	@echo "Stopping services..."
	@docker compose down

.PHONY: logs
logs: ## View logs from all running services
	@docker compose logs -f

.PHONY: ps
ps: ## List running Docker containers for this project
	@docker compose ps

# ==============================================================================
# .NET Commands
# ==============================================================================

.PHONY: build
build: ## Build the .NET application
	@dotnet build

.PHONY: run
run: ## Run the .NET application locally
	@dotnet run --project ./bayt-alhikmah.csproj

.PHONY: test
test: ## Run tests for the .NET application
	@dotnet test

.PHONY: clean
clean: ## Clean build artifacts (bin and obj folders)
	@dotnet clean

# ==============================================================================
# Database Migrations (Entity Framework Core)
# ==============================================================================

.PHONY: migrate-add
migrate-add: ## Create a new database migration. Usage: make migrate-add NAME=InitialCreate
	@dotnet ef migrations add $(NAME)

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations to the database
	@dotnet run --project ./bayt-alhikmah.csproj -- migrate

# ==============================================================================
# Docker Utilities
# ==============================================================================

.PHONY: prune
prune: ## Remove all unused Docker containers, networks, and images
	@docker system prune -a -f