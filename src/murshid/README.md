# Murshid

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

## Development

### Adding New Features

1. **Embeddings**: Add to `app/ml/embeddings.py`
2. **Recommendations**: Add to `app/ml/recommendations.py`
3. **API**: Add endpoints to `app/api/`
4. **Models**: Define schemas in `app/domain/schemas.py`
