# Maktaba Service

Maktaba is the core platform service for Bayt al Hikmah, built with Go 1.25.5.

## Architecture

This service follows Domain-Driven Design (DDD) with bounded contexts:

### Sources Bounded Context
Manages all source types (books, papers, articles, podcasts, videos)
- **Entities**: Source, Author, Tag, Shelf
- **Repository**: PostgreSQL
- **Handler**: REST API endpoints

### Notes Bounded Context
Manages user notes, annotations, and collections
- **Entities**: Note, Profile, Collection, Review, Annotation
- **Repository**: PostgreSQL
- **Handler**: REST API endpoints

## Technology Stack

- **Language**: Go 1.25.5
- **Framework**: Chi (HTTP router)
- **Database**: PostgreSQL 16 with pgx driver
- **Logging**: stdlib log/slog
- **Metrics**: Prometheus client_golang
- **Tracing**: OpenTelemetry

## Project Structure

```
src/maktaba/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── domain/
│   │   ├── sources/        # Sources bounded context
│   │   │   ├── entities/
│   │   │   ├── handler/
│   │   │   ├── service/
│   │   │   └── repository/postgres/
│   │   └── notes/         # Notes bounded context
│   │       ├── entities/
│   │       ├── handler/
│   │       ├── service/
│   │       └── repository/postgres/
│   ├── middleware/         # HTTP middleware
│   └── pkg/              # Shared utilities
├── migrations/            # Database migrations
├── go.mod
└── go.sum
```

## Running Locally

### Prerequisites
- Go 1.25.5
- PostgreSQL 16

### Setup
```bash
cd src/maktaba

# Install dependencies
go mod download

# Set environment variables
export DATABASE_URL="postgres://maktaba:maktaba@localhost:5432/maktaba?sslmode=disable"
export PORT="8080"

# Run migrations
migrate -path migrations -database "$DATABASE_URL" up

# Start server
go run cmd/server/main.go
```

### Building
```bash
go build -o bin/server ./cmd/server
./bin/server
```

## API Endpoints

### Sources
- `POST /sources` - Create source
- `GET /sources` - List sources (supports `limit` and `offset` query params)
- `GET /sources/{id}` - Get source by ID
- `PUT /sources/{id}` - Update source
- `DELETE /sources/{id}` - Delete source

### Notes
- `POST /notes` - Create note
- `GET /notes?user_id={id}` - List notes for user (supports `limit` and `offset`)
- `GET /notes/{id}` - Get note by ID
- `PUT /notes/{id}` - Update note
- `DELETE /notes/{id}` - Delete note

### Health & Metrics
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

## Development

### Adding New Features

1. **Entities**: Add to `internal/domain/{context}/entities/`
2. **Repository**: Implement interface in `internal/domain/{context}/repository/postgres/`
3. **Service**: Add business logic in `internal/domain/{context}/service/`
4. **Handler**: Add HTTP endpoints in `internal/domain/{context}/handler/`
5. **Migrations**: Create migration in `migrations/`

### Code Conventions

- Use `log/slog` for logging (not third-party loggers)
- Use `pgx` driver for PostgreSQL
- Follow Go community standard conventions
- Export types with PascalCase, keep fields JSON-compatible
- Use interfaces for service contracts

## Testing

```bash
go test ./...
```

## Docker

Build and run with Docker Compose:
```bash
make build-maktaba  # Build image
make up              # Start all services
# Or just: docker-compose up -d maktaba
```

## Make Commands

From project root:
- `make build-maktaba` - Build service image
- `make test-maktaba` - Run tests
- `make lint-maktaba` - Lint code
- `make fmt-maktaba` - Format code
- `make deps-maktaba` - Install dependencies
