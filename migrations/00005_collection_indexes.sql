-- +goose Up
CREATE INDEX IF NOT EXISTS idx_collections_source_ids ON collections USING GIN(source_ids);
CREATE INDEX IF NOT EXISTS idx_collections_created_at ON collections(created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_collections_created_at;
DROP INDEX IF EXISTS idx_collections_source_ids;
