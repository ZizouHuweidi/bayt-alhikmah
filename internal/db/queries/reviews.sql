-- name: CreateReview :one
INSERT INTO reviews (id, user_id, source_id, rating, content, is_public)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, source_id, rating, content, is_public, created_at, updated_at;

-- name: GetReviewByID :one
SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
FROM reviews
WHERE id = $1
LIMIT 1;

-- name: ListReviewsByUser :many
SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
FROM reviews
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicReviewsByUser :many
SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
FROM reviews
WHERE user_id = $1 AND is_public = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListReviewsBySource :many
SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
FROM reviews
WHERE source_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicReviewsBySource :many
SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
FROM reviews
WHERE source_id = $1 AND is_public = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateReview :one
UPDATE reviews
SET rating = $2, content = $3, is_public = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, source_id, rating, content, is_public, created_at, updated_at;

-- name: DeleteReview :exec
DELETE FROM reviews WHERE id = $1;
