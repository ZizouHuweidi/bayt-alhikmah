set dotenv-load := true

compose := "podman compose"
database_url := env_var_or_default("DATABASE_URL", "postgres://maktaba:maktaba@localhost:5432/maktaba?sslmode=disable")

default:
    just --list

up:
    {{compose}} up -d

down:
    {{compose}} down

logs:
    {{compose}} logs -f

ps:
    {{compose}} ps

build:
    podman build -f Containerfile -t bayt-alhikmah:latest .

run:
    go run ./cmd/server

test:
    go test ./...

fmt:
    gofmt -w cmd internal pkg

tidy:
    go mod tidy

migrate:
    DATABASE_URL='{{database_url}}' go run ./cmd/migrate up

migrate-down:
    DATABASE_URL='{{database_url}}' go run ./cmd/migrate down

migrate-status:
    DATABASE_URL='{{database_url}}' go run ./cmd/migrate status

migrate-create name:
    go run ./cmd/migrate create {{name}}

seed:
    DATABASE_URL='{{database_url}}' go run ./cmd/seed

db-shell:
    {{compose}} exec postgres psql -U maktaba -d maktaba

health:
    curl -fsS http://localhost:8080/health

frontend-dev:
    cd frontend && bun run dev

frontend-build:
    cd frontend && bun run build

dev: up migrate
