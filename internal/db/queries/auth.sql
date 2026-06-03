-- name: CreateUser :one
INSERT INTO users (id, email, username, password_hash)
VALUES (sqlc.arg(id), LOWER(sqlc.arg(email)), LOWER(sqlc.arg(username)), sqlc.arg(password_hash))
RETURNING id, email, username, password_hash, created_at, updated_at;

-- name: GetUserByEmailOrUsername :one
SELECT id, email, username, password_hash, created_at, updated_at
FROM users
WHERE email = LOWER($1) OR username = LOWER($1)
LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, username, password_hash, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token_hash, family_id, expires_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetRefreshToken :one
SELECT id, user_id, token_hash, family_id, expires_at, revoked_at, created_at
FROM refresh_tokens
WHERE token_hash = $1
LIMIT 1;

-- name: InsertRotatedRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token_hash, family_id, expires_at)
VALUES ($1, $2, $3, $4, $5);

-- name: RevokeRefreshTokenForRotation :execrows
UPDATE refresh_tokens
SET revoked_at = $2, replaced_by_token_id = $3
WHERE id = $1 AND revoked_at IS NULL;

-- name: RevokeRefreshTokenFamily :exec
UPDATE refresh_tokens
SET revoked_at = COALESCE(revoked_at, $2)
WHERE family_id = $1;
