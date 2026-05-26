-- +goose Up
CREATE TABLE IF NOT EXISTS user_library_items (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL CHECK (status IN ('to_consume', 'in_progress', 'completed', 'paused', 'abandoned')),
    progress_value INTEGER CHECK (progress_value IS NULL OR progress_value >= 0),
    progress_unit VARCHAR(50) CHECK (progress_unit IS NULL OR progress_unit IN ('page', 'percent', 'minute', 'second', 'episode')),
    visibility VARCHAR(50) NOT NULL DEFAULT 'private' CHECK (visibility IN ('private', 'unlisted', 'public')),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (user_id, source_id)
);

CREATE INDEX IF NOT EXISTS idx_user_library_items_user_id ON user_library_items(user_id);
CREATE INDEX IF NOT EXISTS idx_user_library_items_source_id ON user_library_items(source_id);
CREATE INDEX IF NOT EXISTS idx_user_library_items_status ON user_library_items(status);
CREATE INDEX IF NOT EXISTS idx_user_library_items_visibility ON user_library_items(visibility);
CREATE INDEX IF NOT EXISTS idx_user_library_items_created_at ON user_library_items(created_at DESC);

CREATE TRIGGER update_user_library_items_updated_at
    BEFORE UPDATE ON user_library_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_user_library_items_updated_at ON user_library_items;
DROP TABLE IF EXISTS user_library_items;
