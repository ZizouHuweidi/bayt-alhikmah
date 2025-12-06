
# Variables
API_PROJECT := src/BaytAlHikmah.Api/BaytAlHikmah.Api.csproj
INFRA_PROJECT := src/BaytAlHikmah.Infrastructure/BaytAlHikmah.Infrastructure.csproj
STARTUP_PROJECT := $(API_PROJECT)
DOTNET_TOOLS := dotnet-ef

.PHONY: help setup build run clean test docker-up docker-down docker-logs migration-add migration-apply db-shell

help:
	@echo "Available commands:"
	@echo "  setup           Initialize project (restore tools, copy .env.example)"
	@echo "  build           Build the solution"
	@echo "  run             Run the API locally"
	@echo "  clean           Clean build artifacts"
	@echo "  test            Run tests"
	@echo "  docker-up       Start services with Docker Compose"
	@echo "  docker-down     Stop services"
	@echo "  docker-logs     View Docker logs"
	@echo "  migration-add   Create a new migration (interactive)"
	@echo "  migration-apply Apply migrations to the database"
	@echo "  db-shell        Connect to the Postgres database shell"

setup:
	@echo "Setting up project..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env from .env.example"; fi
	@dotnet tool restore
	@dotnet restore

build:
	@echo "Building solution..."
	@dotnet build

run:
	@echo "Running API..."
	@dotnet run --project $(API_PROJECT)

clean:
	@echo "Cleaning artifacts..."
	@dotnet clean

test:
	@echo "Running tests..."
	@dotnet test

docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d --build

docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f

migration-add:
	@read -p "Enter migration name: " name; \
	dotnet ef migrations add $$name --project $(INFRA_PROJECT) --startup-project $(STARTUP_PROJECT)

migration-apply:
	@echo "Applying migrations..."
	@dotnet ef database update --project $(INFRA_PROJECT) --startup-project $(STARTUP_PROJECT)

db-shell:
	@echo "Connecting to database..."
	@docker exec -it bayt-alhikmah-postgres-1 psql -U postgres -d bayt_alhikmah
