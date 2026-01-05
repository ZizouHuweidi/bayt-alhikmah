# Murshid Service

Murshid is the ML and recommendations service for Bayt al Hikmah, built with Python 3.11 and FastAPI.

## Responsibilities

- Text embeddings generation
- Content recommendations
- Similarity search
- Vector operations

## Technology Stack

- **Language**: Python 3.13
- **Framework**: FastAPI
- **Dependency Management**: uv
- **ML**: sentence-transformers, PyTorch
- **Vector DB**: pgvector via PostgreSQL
- **Logging**: structlog
- **Metrics**: prometheus_client

## Project Structure

```
src/murshid/
├── app/
│   ├── api/
│   │   └── recommendations.py  # Recommendation endpoints
│   ├── core/
│   │   ├── config.py          # Configuration
│   │   └── logging.py         # Structured logging
│   ├── domain/
│   │   └── schemas.py         # Pydantic models
│   └── ml/
│       ├── embeddings.py       # Text embeddings
│       └── recommendations.py  # Recommendation algorithms
├── main.py
├── pyproject.toml
├── uv.lock
└── Dockerfile
```

## Running Locally

### Prerequisites
- Python 3.13+
- [uv](https://docs.astral.sh/uv/)
- CUDA-capable GPU (optional, for faster embeddings)

### Setup
```bash
cd src/murshid

# Install dependencies and create virtual environment
uv sync

# Set environment variables
export DATABASE_URL="postgres://murshid:murshid@localhost:5432/murshid?sslmode=disable"

# Run server
uv run python main.go
```


## API Endpoints

### Recommendations
- `GET /api/v1/recommendations/{source_id}` - Get similar sources
- `POST /api/v1/embeddings/generate` - Generate embeddings for text

### Health & Metrics
- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe
- `GET /metrics` - Prometheus metrics

## Configuration

Environment variables:

- `DATABASE_URL` - PostgreSQL connection string with pgvector
- `KAFKA_BROKERS` - Kafka bootstrap servers
- `OTEL_EXPORTER_ENDPOINT` - OpenTelemetry collector endpoint
- `OTEL_SERVICE_NAME` - Service name for tracing
- `EMBEDDING_MODEL` - HuggingFace model name (default: sentence-transformers/all-MiniLM-L6-v2)

## ML Models

### Default Model
`sentence-transformers/all-MiniLM-L6-v2`
- Fast inference
- Good quality for semantic search
- Small footprint (80MB)

### Custom Models
Set `EMBEDDING_MODEL` environment variable to use different models.

## Development

### Adding New Features

1. **Embeddings**: Add to `app/ml/embeddings.py`
2. **Recommendations**: Add to `app/ml/recommendations.py`
3. **API**: Add endpoints to `app/api/`
4. **Models**: Define schemas in `app/domain/schemas.py`

### Code Conventions

- Use FastAPI dependency injection
- Follow PEP 8
- Type hints required
- Use structlog for structured logging
- Cache embeddings when possible

## Docker

Build and run with Docker Compose:
```bash
make build-murshid  # Build image
make up              # Start all services
# Or just: docker-compose up -d murshid
```

## Make Commands

From project root:
- `make build-murshid` - Build service image
- `make test-murshid` - Run tests
- `make lint-murshid` - Lint code
- `make fmt-murshid` - Format code
- `make deps-murshid` - Install dependencies
- `make logs-murshid` - View service logs
