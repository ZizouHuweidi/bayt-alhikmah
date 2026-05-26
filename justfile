set dotenv-load := true

compose := "podman compose"

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
    go run ./cmd/migrate up

migrate-down:
    go run ./cmd/migrate down

migrate-status:
    go run ./cmd/migrate status

migrate-create name:
    go run ./cmd/migrate create {{name}}

seed:
    go run ./cmd/seed

db-shell:
    {{compose}} exec postgres psql -U maktaba -d maktaba

health:
    curl -fsS http://localhost:8080/health

frontend-dev:
    cd frontend && bun run dev

frontend-build:
    cd frontend && bun run build

dev: up migrate
