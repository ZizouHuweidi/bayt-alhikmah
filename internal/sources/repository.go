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
