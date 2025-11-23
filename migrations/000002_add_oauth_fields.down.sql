ALTER TABLE users
DROP COLUMN provider,
DROP COLUMN provider_id;

ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
