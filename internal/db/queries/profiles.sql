-- name: GetProfileByUserID :one
SELECT id, user_id, display_name, bio, public_profile, created_at, updated_at
FROM profiles
WHERE user_id = $1
LIMIT 1;

-- name: GetPublicProfileByUsername :one
SELECT p.id, p.user_id, u.username, p.display_name, p.bio, p.public_profile, p.created_at, p.updated_at
FROM profiles p
JOIN users u ON u.id = p.user_id
WHERE u.username = $1 AND p.public_profile = true
LIMIT 1;

-- name: UpsertProfile :one
WITH upserted AS (
    INSERT INTO profiles (id, user_id, display_name, bio, public_profile)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (user_id) DO UPDATE SET
        display_name = EXCLUDED.display_name,
        bio = EXCLUDED.bio,
        public_profile = EXCLUDED.public_profile,
        updated_at = NOW()
    RETURNING id, user_id, display_name, bio, public_profile, created_at, updated_at
)
SELECT p.id, p.user_id, u.username, p.display_name, p.bio, p.public_profile, p.created_at, p.updated_at
FROM upserted p
JOIN users u ON u.id = p.user_id;
