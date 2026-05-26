-- +goose Up
CREATE TABLE IF NOT EXISTS contributors (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS source_contributors (
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    contributor_id UUID NOT NULL REFERENCES contributors(id) ON DELETE CASCADE,
    role VARCHAR(100) NOT NULL DEFAULT 'author',
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (source_id, contributor_id, role)
);

CREATE TABLE IF NOT EXISTS book_metadata (
    source_id UUID PRIMARY KEY REFERENCES sources(id) ON DELETE CASCADE,
    isbn_10 VARCHAR(10),
    isbn_13 VARCHAR(13),
    publisher VARCHAR(255),
    page_count INTEGER CHECK (page_count IS NULL OR page_count >= 0),
    language VARCHAR(32),
    cover_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contributors_name ON contributors(name);
CREATE INDEX IF NOT EXISTS idx_source_contributors_contributor_id ON source_contributors(contributor_id);
CREATE INDEX IF NOT EXISTS idx_source_contributors_role ON source_contributors(role);
CREATE INDEX IF NOT EXISTS idx_book_metadata_isbn_10 ON book_metadata(isbn_10);
CREATE INDEX IF NOT EXISTS idx_book_metadata_isbn_13 ON book_metadata(isbn_13);

CREATE TRIGGER update_contributors_updated_at
    BEFORE UPDATE ON contributors
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_book_metadata_updated_at
    BEFORE UPDATE ON book_metadata
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_book_metadata_updated_at ON book_metadata;
DROP TRIGGER IF EXISTS update_contributors_updated_at ON contributors;
DROP TABLE IF EXISTS book_metadata;
DROP TABLE IF EXISTS source_contributors;
DROP TABLE IF EXISTS contributors;
