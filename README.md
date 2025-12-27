# Bayt al-Hikmah

> بيت الحكمة - House of Wisdom

A modern platform for organizing, engaging with, and tracking knowledge sources across various media.

## 🎯 Current State & Phase

**Current Phase**: Infrastructure Setup & Core Service Development  
**Stage**: Core Infrastructure Complete - Ready for Feature Development

### ✅ Completed

- k3d local development cluster with local registry
- PostgreSQL database with proper credentials (hikmah_db)
- Redpanda event bus (v24.2.5 dev-container mode)
- Grafana + Tempo observability stack  
- ORY Kratos identity service configured
- Maktba (Catalog service) with .NET 10 stable
- EF Core migrations with retry logic at startup
- UUIDv7 for all entity IDs
- KrakenD API Gateway deployed
- Docker containerization with Tiltfile

### 🚧 In Progress

- Gateway routing configuration (debugging 404s)
- Kratos database migrations for full auth
- End-to-end authentication flow

### 🎯 Next Steps

1. Fix Gateway routing to Maktba
2. Run Kratos migrations for auth
3. Implement authentication flow through gateway
4. Create admin dashboard for source management

## 🏗️ Architecture

### Core Services

- **maktba** — Catalog/Source-of-truth (.NET 10) - Sources, authors, taxonomies
- **warraq** — Notes/Profiles/Collections (Go) - User notes, annotations
- **bahith** — Scraper/Ingestion (Go) - External metadata fetching
- **murshid** — Recommendations/ML (Python) - AI-powered recommendations
- **madkhal** — Gateway (KrakenD) - API aggregation and auth

### Infrastructure Stack

- **Database**: PostgreSQL with pgvector
- **Event Bus**: Redpanda (Kafka API)
- **Search**: Meilisearch
- **Auth**: ORY Kratos + Hydra + Oathkeeper
- **Observability**: Grafana stack (Prometheus, Loki, Tempo)
- **Local Dev**: k3d + Tilt

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- k3d
- kubectl
- Tilt
- .NET 10 SDK
- Go 1.21+
- Python 3.11+

### 1. Setup Development Environment

```bash
# Clone and setup
git clone github.com/zizouhuweidi/bayt-alhikmah
cd bayt-alhikmah

# Install dependencies
make setup        # Create k3d cluster
make up          # Launch Tilt dev loop
```

### 2. Verify Services

```bash
# Check cluster status
kubectl get pods

# Port forward for local access
kubectl port-forward svc/maktba 5000:80
kubectl port-forward svc/kratos 4433:4433 4434:4434
kubectl port-forward svc/grafana 3000:3000
```

### 3. Test APIs

```bash
# Create a source
curl -X POST http://localhost:5000/sources \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Book", "type": 0, "description": "A test book", "url": "http://example.com"}'

# Get source
curl http://localhost:5000/sources/{id}

# Health check
curl http://localhost:5000/healthz
```

## 🛠️ Development Commands

```bash
# Cluster Management
make setup      # Setup k3d cluster
make up         # Start development environment
make down       # Stop Tilt and resources
make clean      # Delete entire cluster

# Service Development
make build      # Build all services
make test       # Run tests
make lint       # Run linting

# Database
make db-migrate # Run database migrations
make db-reset   # Reset database

# Utilities
make logs       # Show service logs
make status     # Show cluster status
```

## 📊 Monitoring & Observability

- **Grafana**: http://localhost:3000 (admin/admin)
- **Traces**: Tempo integration
- **Metrics**: Prometheus endpoints
- **Logs**: Structured logging with trace correlation

## 🔧 Configuration

### Environment Variables

Key environment variables are managed through Kubernetes ConfigMaps. See `deploy/overlays/dev/` for development configurations.

### Database Connections

- PostgreSQL: `postgres:5432`
- Connection strings configured per service

### Authentication

- ORY Kratos: Self-hosted identity management
- Admin UI: http://localhost:4434/
- API endpoints: http://localhost:4433/

## 🎯 Development Goals

### Current Sprint Goals

- [x] Fix Maktba service deployment
- [x] Implement source CRUD operations
- [ ] Set up authentication flow (Kratos migrations pending)
- [x] Test end-to-end API functionality

### Next Phase Priorities

- [ ] Fix Gateway routing to backend services
- [ ] Implement Warraq service (notes)
- [ ] Set up Meilisearch integration
- [ ] Create recommendation pipeline
