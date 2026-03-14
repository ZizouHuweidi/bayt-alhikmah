# Bayt al Hikmah

Bayt al Hikmah is a knowledge management platform that helps users organize, explore, discover and interact with knowledge sources.


## Platform Service 
- **Language**: Go
- **Responsibilities**:
  - Sources management (books, papers, articles, etc.)
  - Notes management
  - Collections and reviews
  - User profiles

## Infrastructure

- **PostgreSQL 18**: Primary database with pgvector support
- **Redis 8**: Caching and session storage
- **Redpanda**: Kafka-compatible event streaming
- **Meilisearch v1.5**: Full-text search
- **ORY Stack**: Authentication and authorization (Kratos, Hydra, Oathkeeper)

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.26 (for local development)

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
- `make api` - Open Maktaba API in browser
- `make pgadmin` - Open pgAdmin
- `make help` - See all available commands


## API Endpoints

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
