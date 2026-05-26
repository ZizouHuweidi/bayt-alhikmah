-- +goose Up
DROP INDEX IF EXISTS idx_sources_author_id;
ALTER TABLE sources DROP COLUMN IF EXISTS author_id;
DROP TRIGGER IF EXISTS update_authors_updated_at ON authors;
DROP TABLE IF EXISTS authors;

-- +goose Down
CREATE TABLE IF NOT EXISTS authors (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    birth_date VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE sources ADD COLUMN IF NOT EXISTS author_id UUID REFERENCES authors(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_sources_author_id ON sources(author_id);
CREATE TRIGGER update_authors_updated_at BEFORE UPDATE ON authors FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
