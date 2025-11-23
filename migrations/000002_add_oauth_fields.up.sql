ALTER TABLE users
ADD COLUMN provider VARCHAR(50) DEFAULT '',
ADD COLUMN provider_id VARCHAR(255) DEFAULT '';

-- Make password_hash nullable because OAuth users might not have a password
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;
