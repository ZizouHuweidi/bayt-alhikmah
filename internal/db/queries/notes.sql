-- name: CreateNote :one
INSERT INTO notes (
    user_id, source_id, content, content_type, is_public, annotations, tags
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetNote :one
SELECT * FROM notes WHERE id = $1 LIMIT 1;

-- name: ListNotesByUser :many
SELECT * FROM notes
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListNotesBySource :many
SELECT * FROM notes
WHERE source_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicNotes :many
SELECT * FROM notes
WHERE is_public = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateNote :one
UPDATE notes
SET
    content = COALESCE(sqlc.narg(content), content),
    content_type = COALESCE(sqlc.narg(content_type), content_type),
    is_public = COALESCE(sqlc.narg(is_public), is_public),
    annotations = COALESCE(sqlc.narg(annotations), annotations),
    tags = COALESCE(sqlc.narg(tags), tags),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1;

-- name: CountNotesByUser :one
SELECT COUNT(*) FROM notes WHERE user_id = $1;
