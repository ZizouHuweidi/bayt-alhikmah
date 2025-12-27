# Bayt al-Hikmah

> بيت الحكمة - House of Wisdom

A modern platform for organizing, engaging with, and tracking knowledge sources across various media.

## 🎯 Current State & Phase

**Current Phase**: Phase 3 Complete - Infrastructure & Observability Stabilized  
**Stage**: Identity & Metrics Operational - Feature Implementation Layer Ready

### ✅ Completed

- **Infrastructure**: k3d cluster, PostgreSQL (hikmah_db), Redpanda event bus.
- **Identity**: ORY Kratos configured and migrated (331 migrations applied).
- **Gateway**: KrakenD configured with multi-segment auth routing and query string passthrough.
- **Observability**: Full stack (Grafana, Tempo, Prometheus, Loki) provisioned and verified.
- **Maktba Core**: .NET 10 Catalog service with EF Core auto-migrations and UUIDv7.
- **Telemetry**: Distributed tracing (Tempo), log aggregation (Loki), and metrics scraping (Prometheus).

### 🚧 In Progress

- **Warraq Service**: Notes and annotations (Go implementation).
- **Frontend Integration**: Connecting the Nuxt UI to the KrakenD Gateway.

### 🎯 Next Steps

1. Implement **Warraq** service for user notes.
2. Initialize **Bahith** ingestion pipeline.
3. Finish **Nuxt UI** authentication integration.
4. Setup **Meilisearch** for full-text search across sources.

## 🏗️ Architecture

- **madkhal** (Gateway): KrakenD - Port 8080
- **maktba** (Catalog): .NET 10 - Port 80
- **kratos** (Identity): ORY Kratos - Port 4433 (Public) / 4434 (Admin)
- **Observability**: Grafana (3000), Prometheus (9090), Loki (3100), Tempo (3200)

## 🚀 Quick Start

### 1. Setup Environment
```bash
make setup        # Create k3d cluster
make up           # Launch Tilt dev loop
```

### 2. Verify Auth & APIs (Via Gateway)
```bash
# Register a user
curl -X POST http://localhost:8080/auth/self-service/registration/api

# Create a source
curl -X POST http://localhost:8080/v1/sources \
  -H "Content-Type: application/json" \
  -d '{"title": "The House of Wisdom", "type": 0, "url": "http://example.com"}'
```

## 📊 Monitoring & Observability

- **Grafana**: [http://localhost:3000](http://localhost:3000) (Anonymous Admin)
- **Prometheus**: [http://localhost:9090](http://localhost:9090)
- **Loki Explorer**: Inside Grafana -> Explore -> Select 'Loki'
- **Traces**: Inside Grafana -> Explore -> Select 'Tempo'

## 🛠️ Development Commands

| Command | Description |
|---------|-------------|
| `make setup` | Setup k3d cluster |
| `make up` | Start Tilt dev loop |
| `make port-forward` | Forward all ports for local debug |
| `make test-api` | Run E2E API verification |
| `make db-shell` | Open Postgres shell |
| `make status` | Show cluster status |

## 🎯 Development Goals

- [x] Stabilize PostgreSQL & Redpanda
- [x] Configure Kratos & Apply Migrations
- [x] Refine KrakenD Auth Routing
- [x] Implement Full Observability Stack (Logs/Metrics/Traces)
- [ ] Implement Warraq (Notes) feature
- [ ] Implement Bahith (Ingestion) feature
