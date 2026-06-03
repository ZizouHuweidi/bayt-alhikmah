-- name: CreateSource :one
INSERT INTO sources (id, title, subtitle, type, description, publisher, isbn, doi, url, external_id, tags, published_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id, title, subtitle, type, description, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at;

-- name: InsertBookSource :one
INSERT INTO sources (id, title, subtitle, type, description, publisher, isbn, url, external_id, tags, published_at)
VALUES ($1, $2, $3, 'book', $4, $5, COALESCE($6, $7), $8, $9, $10, $11)
RETURNING id, title, subtitle, type, description, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at;

-- name: InsertBookMetadata :one
INSERT INTO book_metadata (source_id, isbn_10, isbn_13, publisher, page_count, language, cover_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING source_id, isbn_10, isbn_13, publisher, page_count, language, cover_url, created_at, updated_at;

-- name: UpsertContributor :one
INSERT INTO contributors (id, name)
VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING id;

-- name: InsertSourceContributor :one
INSERT INTO source_contributors (source_id, contributor_id, role, position)
VALUES ($1, $2, $3, $4)
RETURNING contributor_id, sqlc.arg('contributor_name')::text AS name, role, position, NOW()::timestamptz AS created_at, NOW()::timestamptz AS updated_at;

-- name: GetSourceByID :one
SELECT id, title, subtitle, type, description, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
FROM sources
WHERE id = $1
LIMIT 1;

-- name: GetBookSourceByID :one
SELECT id, title, subtitle, type, description, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
FROM sources
WHERE id = $1 AND type = 'book'
LIMIT 1;

-- name: GetBookMetadata :one
SELECT source_id, isbn_10, isbn_13, publisher, page_count, language, cover_url, created_at, updated_at
FROM book_metadata
WHERE source_id = $1
LIMIT 1;

-- name: ListContributorsBySource :many
SELECT c.id, c.name, sc.role, sc.position, c.created_at, c.updated_at
FROM source_contributors sc
JOIN contributors c ON c.id = sc.contributor_id
WHERE sc.source_id = $1
ORDER BY sc.position ASC, c.name ASC;

-- name: ListSources :many
SELECT id, title, subtitle, type, description, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
FROM sources
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListSourcesByType :many
SELECT id, title, subtitle, type, description, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
FROM sources
WHERE type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateSource :one
UPDATE sources
SET title = $2, subtitle = $3, type = $4, description = $5, publisher = $6,
    isbn = $7, doi = $8, url = $9, external_id = $10, tags = $11, published_at = $12, updated_at = NOW()
WHERE id = $1
RETURNING id, title, subtitle, type, description, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at;

-- name: DeleteSource :exec
DELETE FROM sources WHERE id = $1;

-- name: SearchSources :many
SELECT id, title, subtitle, type, description, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
FROM sources
WHERE title ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSources :one
SELECT COUNT(*) FROM sources;
