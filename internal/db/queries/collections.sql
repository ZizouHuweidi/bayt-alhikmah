-- name: CreateCollection :one
INSERT INTO collections (id, user_id, name, description, is_public, source_ids)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, name, description, is_public, source_ids, created_at, updated_at;

-- name: GetCollectionByID :one
SELECT id, user_id, name, description, is_public, source_ids, created_at, updated_at
FROM collections
WHERE id = $1
LIMIT 1;

-- name: ListCollectionsByUser :many
SELECT id, user_id, name, description, is_public, source_ids, created_at, updated_at
FROM collections
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicCollectionsByUser :many
SELECT id, user_id, name, description, is_public, source_ids, created_at, updated_at
FROM collections
WHERE user_id = $1 AND is_public = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateCollection :one
UPDATE collections
SET name = $2, description = $3, is_public = $4, source_ids = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, name, description, is_public, source_ids, created_at, updated_at;

-- name: DeleteCollection :exec
DELETE FROM collections WHERE id = $1;

-- name: SourceExists :one
SELECT EXISTS (SELECT 1 FROM sources WHERE id = $1);
