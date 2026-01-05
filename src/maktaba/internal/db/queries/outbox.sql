-- name: CreateOutboxEvent :one
INSERT INTO outbox (
    aggregate_type, aggregate_id, event_type, payload
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUnpublishedEvents :many
SELECT * FROM outbox
WHERE published = false
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkEventPublished :exec
UPDATE outbox
SET published = true, published_at = NOW()
WHERE id = $1;

-- name: DeleteOldPublishedEvents :exec
DELETE FROM outbox
WHERE published = true AND published_at < NOW() - INTERVAL '7 days';
