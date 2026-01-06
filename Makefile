.PHONY: help up down restart logs clean install-deps build-maktaba build-bahith build-murshid build-all test-maktaba test-bahith test-murshid test-all lint-maktaba lint-bahith lint-murshid lint-all migrate-maktaba migrate-all deps-maktaba deps-bahith deps-murshid deps-all

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up:  ## Start all services
	docker-compose up -d

down:  ## Stop all services
	docker-compose down

restart:  ## Restart all services
	docker-compose restart

logs:  ## View logs from all services
	docker-compose logs -f

logs-%:  ## View logs from specific service (e.g., make logs-maktaba)
	docker-compose logs -f $*

clean:  ## Remove all containers, volumes, and networks
	docker-compose down -v
	docker system prune -f

build-maktaba:  ## Build maktaba service image
	cd src/maktaba && docker build -t bayt-alhikmah/maktaba:latest .

build-bahith:  ## Build bahith service image
	cd src/bahith && docker build -t bayt-alhikmah/bahith:latest .

build-murshid:  ## Build murshid service image
	cd src/murshid && docker build -t bayt-alhikmah/murshid:latest .

build-all: build-maktaba build-bahith build-murshid  ## Build all service images

test-maktaba:  ## Run maktaba tests
	cd src/maktaba && go test ./...

test-bahith:  ## Run bahith tests
	cd src/bahith && pytest

test-murshid:  ## Run murshid tests
	cd src/murshid && pytest

test-all: test-maktaba test-bahith test-murshid  ## Run all tests

lint-maktaba:  ## Lint maktaba code
	cd src/maktaba && golangci-lint run

lint-bahith:  ## Lint bahith code
	cd src/bahith && ruff check . && black --check .

lint-murshid:  ## Lint murshid code
	cd src/murshid && ruff check . && black --check .

lint-all: lint-maktaba lint-bahith lint-murshid  ## Lint all code

fmt-maktaba:  ## Format maktaba code
	cd src/maktaba && gofmt -s -w .

fmt-bahith:  ## Format bahith code
	cd src/bahith && ruff check --fix . && black .

fmt-murshid:  ## Format murshid code
	cd src/murshid && ruff check --fix . && black .

fmt-all: fmt-maktaba fmt-bahith fmt-murshid  ## Format all code

migrate-maktaba:  ## Run maktaba database migrations
	docker run --rm -v $(PWD)/src/maktaba/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "postgres://maktaba:maktaba@localhost:5432/maktaba?sslmode=disable" up

migrate-all: migrate-maktaba  ## Run all database migrations

deps-maktaba:  ## Install maktaba dependencies
	cd src/maktaba && go mod download && go mod tidy

deps-bahith:  ## Install bahith dependencies
	cd src/bahith && uv sync

deps-murshid:  ## Install murshid dependencies
	cd src/murshid && uv sync

deps-all: deps-maktaba deps-bahith deps-murshid  ## Install all dependencies

db-shell:  ## Open PostgreSQL shell
	docker-compose exec postgres psql -U postgres

redis-shell:  ## Open Redis CLI
	docker-compose exec redis redis-cli

kafka-shell:  ## Open Redpanda CLI
	docker-compose exec redpanda rpk topic list

grafana:  ## Open Grafana
	@echo "Opening Grafana at http://localhost:3000"
	@xdg-open http://localhost:3000 || open http://localhost:3000 || echo "Please open http://localhost:3000 in your browser"

prometheus:  ## Open Prometheus
	@echo "Opening Prometheus at http://localhost:9090"
	@xdg-open http://localhost:9090 || open http://localhost:9090 || echo "Please open http://localhost:9090 in your browser"

meilisearch:  ## Open Meilisearch
	@echo "Opening Meilisearch at http://localhost:7700"
	@xdg-open http://localhost:7700 || open http://localhost:7700 || echo "Please open http://localhost:7700 in your browser"

pgadmin:  ## Open pgAdmin
	@echo "Opening pgAdmin at http://localhost:5050"
	@xdg-open http://localhost:5050 || open http://localhost:5050 || echo "Please open http://localhost:5050 in your browser"

api:  ## Open Maktaba API
	@echo "Opening Maktaba API at http://localhost:8080"
	@xdg-open http://localhost:8080 || open http://localhost:8080 || echo "Please open http://localhost:8080 in your browser"

ps:  ## Show running containers
	docker-compose ps

status:  ## Show status of all services
	@echo "=== Service Status ==="
	@docker-compose ps
	@echo ""
	@echo "=== Health Checks ==="
	@curl -s http://localhost:8080/health 2>/dev/null && echo "Maktaba: OK" || echo "Maktaba: DOWN"
	@curl -s http://localhost:8003/healthz 2>/dev/null && echo "Bahith: OK" || echo "Bahith: DOWN"
	@curl -s http://localhost:8004/healthz 2>/dev/null && echo "Murshid: OK" || echo "Murshid: DOWN"

seed-db:  ## Seed database with sample data
	docker-compose exec postgres psql -U postgres -d maktaba -c "INSERT INTO sources (title, type) VALUES ('Test Source', 'book');"

clean-volumes:  ## Remove all volumes (WARNING: deletes all data)
	docker-compose down -v
	docker volume rm $$(docker volume ls -q)

rebuild: down build-all up  ## Rebuild and restart all services

sqlc-gen:  ## Generate Go code from SQL queries
	cd src/maktaba && sqlc generate

dev: deps-all up migrate-all  ## Run full local development setup
	@echo "Local development environment is ready!"
	@echo "Maktaba: http://localhost:8080"
	@echo "Grafana: http://localhost:3000"
	@echo "Kratos:  http://localhost:4433"
