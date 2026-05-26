# Bayt al Hikmah

Bayt al Hikmah is a knowledge-management platform for organizing, tracking, reviewing, and annotating knowledge sources.

The current implementation is a Go platform service plus a React frontend. The backend is intentionally simple: one service, PostgreSQL, first-party auth, and explicit SQL through native `pgx`/`pgxpool`.

## Backend

- Go service in `cmd/server`
- Standard library `net/http` router
- PostgreSQL using native `pgx`/`pgxpool`
- Goose SQL migrations in `migrations/`, run through `cmd/migrate`
- Server construction and middleware live in `internal/server`
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

Prerequisites: Go, Node/npm, Podman, Podman Compose, and `just`.

1. Copy environment variables: `cp .env.example .env`
2. Start database and app services: `just up`
3. Run migrations: `just migrate`
4. Seed demo data: `just seed`
5. Start the frontend: `just frontend-dev`
6. Open the frontend at `http://localhost:3000`

Demo login:

- Email: `demo@example.com`
- Password: `password12345`

Local service URLs:

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- Health: `http://localhost:8080/health`
- Readiness: `http://localhost:8080/ready`

Happy-path MVP smoke test:

1. Login with the demo user.
2. Create a book from the dashboard.
3. Add an existing source to the library.
4. Update library status and visibility.
5. Create and delete a note.
6. Create and delete a review.
7. Create and delete a collection.
8. Edit profile settings and enable the public profile.
9. Open `/users/demo_reader/profile`.
10. Open a book detail page from the dashboard.

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
- `just seed` - Seed local demo data (`demo@example.com` / `password12345`)
- `just db-shell` - Open psql inside the Postgres container
- `just health` - Check service liveness endpoint
- `just frontend-dev` - Start the frontend dev server
- `just frontend-build` - Build the frontend

## API Docs

The repo includes a Bruno collection in `bruno/` for local API exploration.

1. Open the `bruno/` directory in Bruno.
2. Select the `local` environment.
3. Run `Auth/Register` or `Auth/Login`.
4. Copy `tokens.access_token` from the response into the `access_token` environment variable.
5. Use the protected `Sources`, `Library`, `Notes`, `Reviews`, and `Collections` requests.

The collection documents the current auth, profile, source, book, library, note, review, and collection endpoints. Refresh uses the `bh_refresh_token` HttpOnly cookie returned by register/login.

## Production Notes

- Set `ENVIRONMENT=production` for production deployments.
- `AUTH_ED25519_PRIVATE_KEY` is required in production and must be a base64 Ed25519 32-byte seed or 64-byte private key.
- If `AUTH_ED25519_PRIVATE_KEY` is omitted in development, the server uses an ephemeral key and existing access tokens become invalid after restart.
- Configure `CORS_ALLOWED_ORIGINS` with the deployed frontend origin.
- `DATABASE_URL` should point to the production PostgreSQL database.

## API Endpoints

Health:

- `GET /health` - Liveness endpoint
- `GET /ready` - Readiness endpoint that verifies database connectivity

Auth:

- `POST /auth/register` - Register with email, username, and password
- `POST /auth/login` - Login with email or username and password
- `POST /auth/refresh` - Rotate refresh token and issue a new access token
- `GET /api/me` - Current authenticated user

Profiles:

- `GET /api/profile` - Get authenticated user's profile, creating an empty profile if needed
- `PUT /api/profile` - Update authenticated user's profile
- `GET /users/{username}/profile` - Get a public profile by username

Sources:

- `POST /api/sources` - Create source
- `POST /api/sources/books` - Create book source with book metadata and contributors
- `GET /sources` - List public sources
- `GET /sources/search?q={query}` - Search sources by title
- `GET /sources/{id}` - Get source by ID
- `GET /sources/books/{id}` - Get book source with book metadata and contributors
- `PUT /api/sources/{id}` - Update source
- `DELETE /api/sources/{id}` - Delete source

Library:

- `POST /api/library/items` - Add a source to the authenticated user's library
- `GET /api/library/items` - List authenticated user's library items
- `GET /api/library/items/with-sources` - List authenticated user's library items with source summaries
- `GET /api/library/items/{id}` - Get own library item by ID
- `PUT /api/library/items/{id}` - Update own library item status/progress/visibility
- `DELETE /api/library/items/{id}` - Remove own library item
- `GET /users/{user}/library` - List public library items for a user by username or user ID
- `GET /users/{user}/library/with-sources` - List public library items with source summaries by username

Collections:

- `POST /api/collections` - Create authenticated user's collection
- `GET /api/collections` - List authenticated user's collections
- `GET /collections?user_id={id}` - List public collections for a user
- `GET /collections/{id}` - Get a public collection by ID
- `PUT /api/collections/{id}` - Update own collection
- `DELETE /api/collections/{id}` - Delete own collection

Notes:

- `POST /api/notes` - Create authenticated user's note
- `GET /api/notes` - List authenticated user's notes
- `GET /notes?public=true` - List public notes
- `GET /notes?source_id={id}` - List public notes for a source
- `GET /notes/{id}` - Get a public note by ID
- `PUT /api/notes/{id}` - Update own note
- `DELETE /api/notes/{id}` - Delete own note

Reviews:

- `POST /api/reviews` - Create authenticated user's source review
- `GET /api/reviews` - List authenticated user's reviews
- `GET /reviews?source_id={id}` - List public reviews for a source
- `GET /reviews?user_id={id}` - List public reviews for a user
- `GET /reviews/{id}` - Get a public review by ID
- `PUT /api/reviews/{id}` - Update own review
- `DELETE /api/reviews/{id}` - Delete own review

Authenticated requests use an EdDSA JWT access token:

```http
Authorization: Bearer <access_token>
```
