# Bayt al Hikmah

Bayt al Hikmah is a knowledge-management platform for organizing, tracking, reviewing, and annotating knowledge sources.

The current implementation is a Go platform service plus a React frontend. The backend is intentionally simple: one service, PostgreSQL, first-party auth, and explicit SQL through native `pgx`/`pgxpool`.

## Backend

- Go service in `cmd/server`
- Standard library `net/http` router
- PostgreSQL using native `pgx`/`pgxpool`
- Goose SQL migrations in `migrations/`, run through `cmd/migrate`
- UUID v7 generated in Go for primary keys
- Argon2id password hashing
- Ed25519/EdDSA JWT access tokens
- Opaque refresh tokens stored hashed in PostgreSQL

## Frontend

- React + Vite
- TanStack Router
- Tailwind CSS

## Infrastructure

- PostgreSQL with pgvector extension available for future recommendation/search work
- Podman-first local development through `compose.yaml` and `justfile`

## Getting Started

1. Copy environment variables: `cp .env.example .env`
2. Start services: `just up`
3. Run migrations: `just migrate`
4. Run tests: `just test`
5. Check health: `just health`

## Useful Commands

- `just up` - Start local services with Podman Compose
- `just down` - Stop local services
- `just logs` - Follow service logs
- `just build` - Build the Go service container image
- `just run` - Run the Go service locally
- `just test` - Run Go tests
- `just fmt` - Format Go code
- `just migrate` - Apply database migration
- `just migrate-down` - Roll back database migration
- `just migrate-status` - Show Goose migration status
- `just migrate-create name` - Create a new Goose SQL migration
- `just db-shell` - Open psql inside the Postgres container
- `just frontend-dev` - Start the frontend dev server

## API Endpoints

Auth:

- `POST /auth/register` - Register with email, username, and password
- `POST /auth/login` - Login with email or username and password
- `POST /auth/refresh` - Rotate refresh token and issue a new access token
- `GET /api/me` - Current authenticated user

Sources:

- `POST /api/sources` - Create source
- `GET /sources` - List public sources
- `GET /sources/search?q={query}` - Search sources by title
- `GET /sources/{id}` - Get source by ID
- `PUT /api/sources/{id}` - Update source
- `DELETE /api/sources/{id}` - Delete source

Notes:

- `POST /api/notes` - Create authenticated user's note
- `GET /notes?public=true` - List public notes
- `GET /notes?source_id={id}` - List public notes for a source
- `GET /notes/{id}` - Get a public note by ID
- `PUT /api/notes/{id}` - Update own note
- `DELETE /api/notes/{id}` - Delete own note

Authenticated requests use an EdDSA JWT access token:

```http
Authorization: Bearer <access_token>
```
