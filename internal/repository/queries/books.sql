-- name: ListBooks :many
SELECT *
FROM books
WHERE (
        $1::text = ''
        OR title ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC;
-- name: GetBook :one
SELECT *
FROM books
WHERE id = $1
LIMIT 1;
-- name: CreateBook :one
INSERT INTO books (
        title,
        author,
        description,
        thumbnail_url
    )
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: UpdateBook :one
UPDATE books
SET title = COALESCE($2, title),
    author = COALESCE($3, author),
    description = COALESCE($4, description),
    thumbnail_url = COALESCE($5, thumbnail_url),
    updated_at = NOW()
WHERE id = $1
RETURNING *;
-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;