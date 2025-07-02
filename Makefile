# Makefile for Bayt al-Hikmah

.PHONY: all build run test clean docker-up docker-down docker-build migrate add-migration remove-migration db-drop db-reset

# Variables
API_PROJECT = src/BaytAlHikmah.Api/BaytAlHikmah.Api.csproj
INFRASTRUCTURE_PROJECT = src/BaytAlHikmah.Infrastructure/BaytAlHikmah.Infrastructure.csproj

# Default target
all: build

# Build the solution
build:
	@echo "Building the solution..."
	dotnet build

# Run the API locally
run:
	@echo "Running the API..."
	dotnet run --project $(API_PROJECT)

# Run the tests
test:
	@echo "Running tests..."
	dotnet test

# Clean the solution
clean:
	@echo "Cleaning the solution..."
	dotnet clean

# Start all services with Docker Compose
docker-up:
	@echo "Starting all services with Docker Compose..."
	docker-compose up -d

# Stop all services with Docker Compose
docker-down:
	@echo "Stopping all services with Docker Compose..."
	docker-compose down

# Build the Docker images
docker-build:
	@echo "Building Docker images..."
	docker-compose build

# Apply pending database migrations
migrate:
	@echo "Applying database migrations..."
	dotnet ef database update --project $(INFRASTRUCTURE_PROJECT)

# Add a new migration (Usage: make add-migration n=MigrationName)
add-migration:
ifndef n
	$(error n is undefined. Usage: make add-migration n=MigrationName)
endif
	@echo "Adding new migration: $(n)..."
	dotnet ef migrations add $(n) --project $(INFRASTRUCTURE_PROJECT)

# Remove the last migration
remove-migration:
	@echo "Removing last migration..."
	dotnet ef migrations remove --project $(INFRASTRUCTURE_PROJECT)

# Drop the database
db-drop:
	@echo "Dropping the database..."
	dotnet ef database drop --project $(INFRASTRUCTURE_PROJECT)

# Reset the database (drop and re-create)
db-reset: db-drop migrate
	@echo "Database reset complete."