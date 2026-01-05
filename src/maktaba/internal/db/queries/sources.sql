-- name: CreateSource :one
INSERT INTO sources (
    title, subtitle, type, description, author_id, publisher,
    isbn, doi, url, external_id, tags, published_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetSource :one
SELECT * FROM sources WHERE id = $1 LIMIT 1;

-- name: ListSources :many
SELECT * FROM sources
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListSourcesByType :many
SELECT * FROM sources
WHERE type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateSource :one
UPDATE sources
SET
    title = COALESCE(sqlc.narg(title), title),
    subtitle = COALESCE(sqlc.narg(subtitle), subtitle),
    type = COALESCE(sqlc.narg(type), type),
    description = COALESCE(sqlc.narg(description), description),
    author_id = COALESCE(sqlc.narg(author_id), author_id),
    publisher = COALESCE(sqlc.narg(publisher), publisher),
    isbn = COALESCE(sqlc.narg(isbn), isbn),
    doi = COALESCE(sqlc.narg(doi), doi),
    url = COALESCE(sqlc.narg(url), url),
    external_id = COALESCE(sqlc.narg(external_id), external_id),
    tags = COALESCE(sqlc.narg(tags), tags),
    published_at = COALESCE(sqlc.narg(published_at), published_at),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteSource :exec
DELETE FROM sources WHERE id = $1;

-- name: CountSources :one
SELECT COUNT(*) FROM sources;

-- name: SearchSourcesByTitle :many
SELECT * FROM sources
WHERE title ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
