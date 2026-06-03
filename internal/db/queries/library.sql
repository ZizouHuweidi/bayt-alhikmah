-- name: CreateLibraryItem :one
INSERT INTO user_library_items (id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at;

-- name: GetLibraryItemByID :one
SELECT id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at
FROM user_library_items
WHERE id = $1
LIMIT 1;

-- name: ListLibraryItemsByUser :many
SELECT id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at
FROM user_library_items
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicLibraryItemsByUser :many
SELECT id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at
FROM user_library_items
WHERE user_id = $1 AND visibility = 'public'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicLibraryItemsByUsername :many
SELECT uli.id, uli.user_id, uli.source_id, uli.status, uli.progress_value, uli.progress_unit, uli.visibility, uli.started_at, uli.completed_at, uli.created_at, uli.updated_at
FROM user_library_items uli
JOIN users u ON u.id = uli.user_id
WHERE u.username = $1 AND uli.visibility = 'public'
ORDER BY uli.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLibraryItemsByUserWithSources :many
SELECT uli.id, uli.user_id, uli.source_id, uli.status, uli.progress_value, uli.progress_unit, uli.visibility, uli.started_at, uli.completed_at, uli.created_at, uli.updated_at,
       s.title, s.subtitle, s.type, s.publisher, s.isbn
FROM user_library_items uli
JOIN sources s ON s.id = uli.source_id
WHERE uli.user_id = $1
ORDER BY uli.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicLibraryItemsByUsernameWithSources :many
SELECT uli.id, uli.user_id, uli.source_id, uli.status, uli.progress_value, uli.progress_unit, uli.visibility, uli.started_at, uli.completed_at, uli.created_at, uli.updated_at,
       s.title, s.subtitle, s.type, s.publisher, s.isbn
FROM user_library_items uli
JOIN users u ON u.id = uli.user_id
JOIN sources s ON s.id = uli.source_id
WHERE u.username = $1 AND uli.visibility = 'public'
ORDER BY uli.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateLibraryItem :one
UPDATE user_library_items
SET status = $2, progress_value = $3, progress_unit = $4, visibility = $5, started_at = $6, completed_at = $7, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at;

-- name: DeleteLibraryItem :exec
DELETE FROM user_library_items WHERE id = $1;
