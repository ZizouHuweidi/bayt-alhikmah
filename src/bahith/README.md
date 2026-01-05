# Bahith Service

Bahith is the scraper service for Bayt al Hikmah, built with Python 3.11 and FastAPI.

## Responsibilities

- Web scraping for Islamic sources
- Content ingestion and normalization
- Event publishing to Kafka
- Metadata extraction

## Technology Stack

- **Language**: Python 3.13
- **Framework**: FastAPI
- **Dependency Management**: uv
- **Scraping**: BeautifulSoup4, Scrapy
- **Event Streaming**: aiokafka
- **Logging**: structlog
- **Metrics**: prometheus_client

## Project Structure

```
src/bahith/
├── app/
│   ├── api/
│   │   └── scrape.py      # Scraping endpoints
│   ├── core/
│   │   ├── config.py      # Configuration
│   │   ├── logging.py     # Structured logging
│   │   └── kafka.py       # Kafka producer
│   ├── domain/
│   │   └── schemas.py     # Pydantic models
│   └── scraper/
│       ├── base.py        # Base scraper
│       └── ...
├── main.py
├── pyproject.toml
├── uv.lock
└── Dockerfile
```

## Running Locally

### Prerequisites
- Python 3.13+
- [uv](https://docs.astral.sh/uv/)

### Setup
```bash
cd src/bahith

# Install dependencies and create virtual environment
uv sync

# Set environment variables
export KAFKA_BROKERS="localhost:9092"

# Run server
uv run python main.py
```


## API Endpoints

### Scraping
- `POST /api/v1/scrape` - Trigger scraping for a URL
- `GET /api/v1/scrape/status/{id}` - Get scrape status

### Health & Metrics
- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe
- `GET /metrics` - Prometheus metrics

## Events

Published to Kafka topics:

- `source-created` - When a new source is scraped
- `source-updated` - When a source is updated
- `scrape-error` - When scraping fails

## Configuration

Environment variables:

- `DATABASE_URL` - SQLite database path (default: sqlite:///./bahith.db)
- `KAFKA_BROKERS` - Kafka bootstrap servers (default: localhost:9092)
- `OTEL_EXPORTER_ENDPOINT` - OpenTelemetry collector endpoint
- `OTEL_SERVICE_NAME` - Service name for tracing

## Development

### Adding New Scrapers

1. Extend `base.Scraper` class
2. Implement `scrape()` method
3. Register in API endpoints
4. Add event publishing logic

### Code Conventions

- Use FastAPI dependency injection
- Follow PEP 8
- Type hints required
- Use structlog for structured logging
- Publish events for all significant actions

## Docker

Build and run with Docker Compose:
```bash
make build-bahith  # Build image
make up            # Start all services
# Or just: docker-compose up -d bahith
```

## Make Commands

From project root:
- `make build-bahith` - Build service image
- `make test-bahith` - Run tests
- `make lint-bahith` - Lint code
- `make fmt-bahith` - Format code
- `make deps-bahith` - Install dependencies
- `make logs-bahith` - View service logs
