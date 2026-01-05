-- Drop tables in reverse order of dependencies
DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
DROP TRIGGER IF EXISTS update_collections_updated_at ON collections;
DROP TRIGGER IF EXISTS update_notes_updated_at ON notes;
DROP TRIGGER IF EXISTS update_profiles_updated_at ON profiles;
DROP TRIGGER IF EXISTS update_sources_updated_at ON sources;
DROP TRIGGER IF EXISTS update_tags_updated_at ON tags;
DROP TRIGGER IF EXISTS update_authors_updated_at ON authors;

DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS sources;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS authors;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop extensions (optional, comment out if you want to keep them)
-- DROP EXTENSION IF EXISTS vector;
-- DROP EXTENSION IF EXISTS "uuid-ossp";
