include .env
export

.PHONY: init
init:
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@go install github.com/segmentio/golines@latest
	@go install github.com/jesseduffield/lazygit@latest
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/mfridman/tparse@latest

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# =====================================================================
# Migrations using golang-migrate via Docker
# =====================================================================
/PHONY: migrate/up migrate/up/all migrate/down migrate/down/all migrate/force
## migrate/up n=<number>: migrates up n steps
migrate/up:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up $(n)

## migrate/up/all: migrates up to latest
migrate/up/all:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up

## migrate/down n=<number>: migrates down n steps
migrate/down:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down $(n)

## migrate/down/all: migrates down all steps
migrate/down/all:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down -all
## migration n=<file_name>: creates migration files up/down for file_name
migration:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ create -seq -ext=.sql -dir=./migrations $(n)
## migrate/force n=<version>: forces migration version number
migrate/force:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database=$(CONNECTION_STRING) force $(n)

refresh: migrate.down.all migrate.up seed


# =====================================================================
# Build
# =====================================================================

# Build the Go binary (backend)
build:
	@echo "Building Go backend..."
	@go build -o main cmd/api/main.go

# Alternative Go builds
build/local:
	@echo "Building with 'local' tags..."
	@go build -tags local -o main .

# Run targets:
# run-backend: uses air (live reload) for backend
run-backend:
	@echo "Running backend (using air)..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		read -p "Go's 'air' is not installed. Install it now? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/cosmtrek/air@latest; \
			air; \
		else \
			echo "air not installed. Exiting..."; exit 1; \
		fi; \
	fi

# run-frontend: install and run client dev server
run-frontend:
	@echo "Running frontend (React dev)..."
	@npm install --prefer-offline --no-fund --prefix ./client
	@npm run dev --prefix ./client

# run: run both backend and frontend concurrently
run:
	@$(MAKE) run-backend & \
	$(MAKE) run-frontend

# =====================================================================
# Testing
# =====================================================================
test:
	@echo "Running unit tests..."
	@go test ./... -v

test-race:
	@echo "Running tests with race detection..."
	@go test -race ./... -v

itest:
	@echo "Running integration tests for internal/database..."
	@go test ./internal/database -v

## test: use hurl to run tests on a running application
# test:
# 	hurl --variables-file .env --file-root . --test .

# =====================================================================
# Docker Compose Operations
# =====================================================================
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# =====================================================================
# Database and Docker Utilities
# =====================================================================
.PHONY: db/conn prune ps inspect prune-dangled-volumes

db/conn:
	psql $(CONNECTION_STRING)

prune:
	docker system prune -a -f --volumes

prune-dangled-volumes:
	docker volume ls -q -f dangling=true | xargs -r docker volume rm

ps:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.RunningFor}}\t{{.Size}}\t{{.Ports}}"

inspect:
	docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" $(n)

# =====================================================================
# Go Module Management
# =====================================================================
.PHONY: list update
list:
	go list -m -u

update:
	go get -u ./...

# =====================================================================
# Quality Control: Lint, Format, and Test
# =====================================================================
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	gofumpt -l -w -extra .
	golines -w .
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# =====================================================================
# psql Utility
# =====================================================================
.PHONY: psql
psql:
	docker run -it --rm --network ${NETWORK} albayt psql -h ${DB_HOST} -U ${DB_USERNAME}
