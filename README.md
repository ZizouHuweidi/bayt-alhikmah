# Bayt al Hikmah

Bayt al Hikmah is an Islamic knowledge management platform that helps users organize, explore, and discover Islamic texts and scholarly resources.

## Services

### Platform Service (Maktaba)
- **Port**: 8080
- **Language**: Go
- **Responsibilities**:
  - Sources management (books, papers, articles, etc.)
  - Notes management
  - Collections and reviews
  - User profiles

### Scraper Service (Bahith)
- **Port**: 8003
- **Language**: Python (FastAPI)
- **Responsibilities**:
  - Web scraping for Islamic sources
  - Content ingestion
  - Event publishing to Kafka

### ML Service (Murshid)
- **Port**: 8004
- **Language**: Python (FastAPI)
- **Responsibilities**:
  - Text embeddings
  - Content recommendations
  - Similarity search

## Infrastructure

- **PostgreSQL 18**: Primary database with pgvector support
- **Redis 8**: Caching and session storage
- **Redpanda**: Kafka-compatible event streaming
- **Meilisearch v1.5**: Full-text search
- **ORY Stack**: Authentication and authorization (Kratos, Hydra, Oathkeeper)

### Observability Stack

- **Prometheus**: Metrics collection
- **Grafana**: Metrics visualization
- **Loki**: Log aggregation
- **Tempo**: Distributed tracing
- **Alertmanager**: Alert management

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.25.5 (for local development)
- Python 3.13+ (for local development)
- [uv](https://docs.astral.sh/uv/) (for Python services)

### Using Docker Compose

1. Clone the repository
2. Copy environment variables: `cp .env.example .env`
3. Start all services: `make up` (or `docker-compose up -d`)
4. Run migrations: `make migrate-maktaba`
5. Check service health: `make status`

#### Useful Make Commands
- `make up` - Start all services
- `make down` - Stop all services
- `make logs` - View all logs
- `make status` - Check service health
- `make build-all` - Build all service images
- `make migrate-all` - Run all migrations
- `make test-all` - Run all tests
- `make lint-all` - Lint all code
- `make clean` - Remove all containers and volumes
- `make api` - Open API in browser
- `make grafana` - Open Grafana dashboard
- `make help` - See all available commands

### Local Development

#### Maktaba (Go)
```bash
cd src/maktaba
go mod download
go run cmd/server/main.go
# Or use make: make deps-maktaba && make build-maktaba
```

#### Bahith (Python)
```bash
cd src/bahith
uv sync
uv run python main.py
```

#### Murshid (Python)
```bash
cd src/murshid
uv sync
uv run python main.py
```


## API Endpoints

### Maktaba

**Sources**:
- `POST /sources` - Create source
- `GET /sources` - List sources
- `GET /sources/{id}` - Get source by ID
- `PUT /sources/{id}` - Update source
- `DELETE /sources/{id}` - Delete source

**Notes**:
- `POST /notes` - Create note
- `GET /notes?user_id={id}` - List notes for user
- `GET /notes/{id}` - Get note by ID
- `PUT /notes/{id}` - Update note
- `DELETE /notes/{id}` - Delete note

### Observability

- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics
