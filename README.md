# Bayt al Hikmah

Bayt al Hikmah is a personal knowledge library for collecting, tracking, annotating, and reviewing sources such as books and future learning media.

The MVP focuses on helping a reader build a library, track reading progress, write notes and reviews, organize sources into collections, and publish a lightweight public profile.

## Core Features

- First-party user registration and login.
- Book/source creation with metadata and contributors.
- Personal library tracking with status, progress, and visibility.
- Notes, reviews, and collections for organizing learning.
- Public profiles and public library views.
- Demo seed data for local testing.
- Bruno API collection for local API exploration.

## Stack

- Backend: Go, Echo, PostgreSQL, pgx/pgxpool.
- Database: Goose migrations, sqlc-generated query code, UUID v7 IDs.
- Auth: Argon2id password hashes, Ed25519 JWT access tokens, opaque refresh tokens.
- Frontend: React Router SPA, React Query, Zustand, React Hook Form, Zod, Tailwind/shadcn-style UI, Biome.
- Tooling: Podman Compose, Containerfiles, justfile commands.

## Local Development

Prerequisites: Go, Node/npm, Podman, Podman Compose, and `just`.

```sh
cp .env.example .env
just up
just migrate
just seed
just frontend-dev
```

Open the frontend at `http://localhost:3000`.

Demo login:

- Email: `demo@example.com`
- Password: `password12345`

Useful URLs:

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- Health: `http://localhost:8080/health`
- Readiness: `http://localhost:8080/ready`

## Common Commands

- `just up` - Start local services.
- `just down` - Stop local services.
- `just migrate` - Apply database migrations.
- `just migrate-status` - Show migration status.
- `just migrate-create name` - Create a migration file.
- `just seed` - Seed demo data.
- `just test` - Run Go tests.
- `just check` - Run sqlc checks, Go tests, and frontend checks.
- `just frontend-dev` - Start the frontend dev server.
- `just frontend-build` - Build the frontend.

## Database Workflow

- Schema changes live in `migrations/` and are applied with Goose through `cmd/migrate`.
- Application queries live in `internal/db/queries/` and generate Go code with sqlc.
- sqlc reads the migrations as schema input, but it does not apply migrations.
- sqlc runs through a container, so a local `sqlc` binary is not required.
- Goose is used as a Go library through `go run ./cmd/migrate`, so a local `goose` binary is not required.

## API Exploration

Open the `bruno/` directory in Bruno and select the `local` environment. The collection covers auth, profile, sources, books, library, notes, reviews, and collections.
