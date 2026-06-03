-- name: CreateNote :one
INSERT INTO notes (id, user_id, source_id, content, content_type, is_public, annotations, tags)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetNoteByID :one
SELECT *
FROM notes
WHERE id = $1
LIMIT 1;

-- name: ListNotesByUser :many
SELECT *
FROM notes
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicNotesByUser :many
SELECT *
FROM notes
WHERE user_id = $1 AND is_public = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListNotesBySource :many
SELECT *
FROM notes
WHERE source_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicNotesBySource :many
SELECT *
FROM notes
WHERE source_id = $1 AND is_public = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicNotes :many
SELECT *
FROM notes
WHERE is_public = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountNotesByUser :one
SELECT COUNT(*) FROM notes WHERE user_id = $1;

-- name: UpdateNote :one
UPDATE notes
SET source_id = $2, content = $3, content_type = $4, is_public = $5, annotations = $6, tags = $7, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1;
