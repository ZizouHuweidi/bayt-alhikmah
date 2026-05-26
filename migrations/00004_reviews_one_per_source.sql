-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_user_source_unique ON reviews(user_id, source_id);
CREATE INDEX IF NOT EXISTS idx_reviews_is_public ON reviews(is_public);
CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_reviews_created_at;
DROP INDEX IF EXISTS idx_reviews_is_public;
DROP INDEX IF EXISTS idx_reviews_user_source_unique;
