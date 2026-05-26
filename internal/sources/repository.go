package sources

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zizouhuweidi/maktaba/internal/db"
)

type postgresRepository struct {
	db *db.DB
}

type sourceRow struct {
	ID          pgtype.UUID
	Title       string
	Subtitle    pgtype.Text
	Type        string
	Description pgtype.Text
	AuthorID    pgtype.UUID
	Publisher   pgtype.Text
	ISBN        pgtype.Text
	DOI         pgtype.Text
	URL         pgtype.Text
	ExternalID  pgtype.Text
	Tags        []byte
	PublishedAt pgtype.Timestamptz
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type bookMetadataRow struct {
	SourceID  pgtype.UUID
	ISBN10    pgtype.Text
	ISBN13    pgtype.Text
	Publisher pgtype.Text
	PageCount pgtype.Int4
	Language  pgtype.Text
	CoverURL  pgtype.Text
	CreatedAt time.Time
	UpdatedAt time.Time
}

type contributorRow struct {
	ID        pgtype.UUID
	Name      string
	Role      string
	Position  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{db: d}
}

func (r *postgresRepository) Create(ctx context.Context, s *Source) (*Source, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	tags, err := json.Marshal(s.Tags)
	if err != nil {
		return nil, err
	}

	row, err := scanSource(r.db.QueryRow(ctx, `
		INSERT INTO sources (id, title, subtitle, type, description, author_id, publisher, isbn, doi, url, external_id, tags, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, title, subtitle, type, description, author_id, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
	`, id.String(), s.Title, s.Subtitle, string(s.Type), s.Description, nullableUUIDString(s.AuthorID), s.Publisher, s.ISBN, s.DOI, s.URL, s.ExternalID, tags, s.PublishedAt))
	if err != nil {
		return nil, err
	}

	return r.mapRow(row), nil
}

func (r *postgresRepository) CreateBook(ctx context.Context, params CreateBookParams) (*Book, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	sourceID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	tags, err := json.Marshal(params.Tags)
	if err != nil {
		return nil, err
	}

	sourceRow, err := scanSource(tx.QueryRow(ctx, `
		INSERT INTO sources (id, title, subtitle, type, description, publisher, isbn, url, external_id, tags, published_at)
		VALUES ($1, $2, $3, 'book', $4, $5, COALESCE($6, $7), $8, $9, $10, $11)
		RETURNING id, title, subtitle, type, description, author_id, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
	`, sourceID.String(), params.Title, params.Subtitle, params.Description, params.Publisher, params.ISBN13, params.ISBN10, params.URL, params.ExternalID, tags, params.PublishedAt))
	if err != nil {
		return nil, err
	}

	metadataRow, err := scanBookMetadata(tx.QueryRow(ctx, `
		INSERT INTO book_metadata (source_id, isbn_10, isbn_13, publisher, page_count, language, cover_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING source_id, isbn_10, isbn_13, publisher, page_count, language, cover_url, created_at, updated_at
	`, sourceID.String(), params.ISBN10, params.ISBN13, params.Publisher, params.PageCount, params.Language, params.CoverURL))
	if err != nil {
		return nil, err
	}

	contributors := make([]*Contributor, 0, len(params.Contributors))
	for position, contributor := range params.Contributors {
		role := contributor.Role
		if role == "" {
			role = "author"
		}

		contributorID, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		var existingID pgtype.UUID
		err = tx.QueryRow(ctx, `
			INSERT INTO contributors (id, name)
			VALUES ($1, $2)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, contributorID.String(), contributor.Name).Scan(&existingID)
		if err != nil {
			return nil, err
		}

		row, err := scanContributor(tx.QueryRow(ctx, `
			INSERT INTO source_contributors (source_id, contributor_id, role, position)
			VALUES ($1, $2, $3, $4)
			RETURNING contributor_id, $5::text, role, position, NOW(), NOW()
		`, sourceID.String(), uuid.UUID(existingID.Bytes).String(), role, position, contributor.Name))
		if err != nil {
			return nil, err
		}
		contributors = append(contributors, mapContributor(row))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &Book{Source: r.mapRow(sourceRow), Metadata: mapBookMetadata(metadataRow), Contributors: contributors}, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Source, error) {
	row, err := scanSource(r.db.QueryRow(ctx, `
		SELECT id, title, subtitle, type, description, author_id, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
		FROM sources WHERE id = $1 LIMIT 1
	`, id.String()))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.mapRow(row), nil
}

func (r *postgresRepository) GetBookByID(ctx context.Context, id uuid.UUID) (*Book, error) {
	sourceRow, err := scanSource(r.db.QueryRow(ctx, `
		SELECT id, title, subtitle, type, description, author_id, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
		FROM sources WHERE id = $1 AND type = 'book' LIMIT 1
	`, id.String()))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	metadataRow, err := scanBookMetadata(r.db.QueryRow(ctx, `
		SELECT source_id, isbn_10, isbn_13, publisher, page_count, language, cover_url, created_at, updated_at
		FROM book_metadata WHERE source_id = $1 LIMIT 1
	`, id.String()))
	if errors.Is(err, pgx.ErrNoRows) {
		return &Book{Source: r.mapRow(sourceRow)}, nil
	}
	if err != nil {
		return nil, err
	}

	contributors, err := r.listContributors(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Book{Source: r.mapRow(sourceRow), Metadata: mapBookMetadata(metadataRow), Contributors: contributors}, nil
}

func (r *postgresRepository) List(ctx context.Context, limit, offset int) ([]*Source, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, subtitle, type, description, author_id, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
		FROM sources ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *postgresRepository) ListByType(ctx context.Context, sourceType SourceType, limit, offset int) ([]*Source, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, subtitle, type, description, author_id, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
		FROM sources WHERE type = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, string(sourceType), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *postgresRepository) Update(ctx context.Context, s *Source) (*Source, error) {
	tags, err := json.Marshal(s.Tags)
	if err != nil {
		return nil, err
	}

	row, err := scanSource(r.db.QueryRow(ctx, `
		UPDATE sources
		SET title = $2, subtitle = $3, type = $4, description = $5, author_id = $6, publisher = $7,
			isbn = $8, doi = $9, url = $10, external_id = $11, tags = $12, published_at = $13, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, subtitle, type, description, author_id, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
	`, s.ID.String(), s.Title, s.Subtitle, string(s.Type), s.Description, nullableUUIDString(s.AuthorID), s.Publisher, s.ISBN, s.DOI, s.URL, s.ExternalID, tags, s.PublishedAt))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.mapRow(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sources WHERE id = $1`, id.String())
	return err
}

func (r *postgresRepository) Search(ctx context.Context, query string, limit, offset int) ([]*Source, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, subtitle, type, description, author_id, publisher, isbn, doi, url, external_id, tags, published_at, created_at, updated_at
		FROM sources WHERE title ILIKE '%' || $1 || '%' ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *postgresRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM sources`).Scan(&count)
	return count, err
}

func (r *postgresRepository) listContributors(ctx context.Context, sourceID uuid.UUID) ([]*Contributor, error) {
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.name, sc.role, sc.position, c.created_at, c.updated_at
		FROM source_contributors sc
		JOIN contributors c ON c.id = sc.contributor_id
		WHERE sc.source_id = $1
		ORDER BY sc.position ASC, c.name ASC
	`, sourceID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contributors := []*Contributor{}
	for rows.Next() {
		row, err := scanContributor(rows)
		if err != nil {
			return nil, err
		}
		contributors = append(contributors, mapContributor(row))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contributors, nil
}

func (r *postgresRepository) scanRows(rows pgx.Rows) ([]*Source, error) {
	var sources []*Source
	for rows.Next() {
		row, err := scanSource(rows)
		if err != nil {
			return nil, err
		}
		sources = append(sources, r.mapRow(row))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sources, nil
}

func scanSource(row pgx.Row) (sourceRow, error) {
	var s sourceRow
	err := row.Scan(
		&s.ID,
		&s.Title,
		&s.Subtitle,
		&s.Type,
		&s.Description,
		&s.AuthorID,
		&s.Publisher,
		&s.ISBN,
		&s.DOI,
		&s.URL,
		&s.ExternalID,
		&s.Tags,
		&s.PublishedAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	return s, err
}

func scanBookMetadata(row pgx.Row) (bookMetadataRow, error) {
	var metadata bookMetadataRow
	err := row.Scan(
		&metadata.SourceID,
		&metadata.ISBN10,
		&metadata.ISBN13,
		&metadata.Publisher,
		&metadata.PageCount,
		&metadata.Language,
		&metadata.CoverURL,
		&metadata.CreatedAt,
		&metadata.UpdatedAt,
	)
	return metadata, err
}

func scanContributor(row pgx.Row) (contributorRow, error) {
	var contributor contributorRow
	err := row.Scan(
		&contributor.ID,
		&contributor.Name,
		&contributor.Role,
		&contributor.Position,
		&contributor.CreatedAt,
		&contributor.UpdatedAt,
	)
	return contributor, err
}

func (r *postgresRepository) mapRow(row sourceRow) *Source {
	var tags []string
	if len(row.Tags) > 0 {
		_ = json.Unmarshal(row.Tags, &tags)
	}

	return &Source{
		ID:          uuid.UUID(row.ID.Bytes),
		Title:       row.Title,
		Subtitle:    stringPtr(row.Subtitle),
		Type:        SourceType(row.Type),
		Description: stringPtr(row.Description),
		AuthorID:    uuidPtr(row.AuthorID),
		Publisher:   stringPtr(row.Publisher),
		ISBN:        stringPtr(row.ISBN),
		DOI:         stringPtr(row.DOI),
		URL:         stringPtr(row.URL),
		ExternalID:  stringPtr(row.ExternalID),
		Tags:        tags,
		PublishedAt: timePtr(row.PublishedAt),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func stringPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func intPtr(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	result := int(value.Int32)
	return &result
}

func mapBookMetadata(row bookMetadataRow) *BookMetadata {
	return &BookMetadata{
		SourceID:  uuid.UUID(row.SourceID.Bytes),
		ISBN10:    stringPtr(row.ISBN10),
		ISBN13:    stringPtr(row.ISBN13),
		Publisher: stringPtr(row.Publisher),
		PageCount: intPtr(row.PageCount),
		Language:  stringPtr(row.Language),
		CoverURL:  stringPtr(row.CoverURL),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func mapContributor(row contributorRow) *Contributor {
	return &Contributor{
		ID:        uuid.UUID(row.ID.Bytes),
		Name:      row.Name,
		Role:      row.Role,
		Position:  row.Position,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func uuidPtr(value pgtype.UUID) *uuid.UUID {
	if !value.Valid {
		return nil
	}
	id := uuid.UUID(value.Bytes)
	return &id
}

func nullableUUIDString(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
}

func timePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}
