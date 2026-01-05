package sources

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zizouhuweidi/maktaba/internal/db"
)

type postgresRepository struct {
	db *db.DB
}

// NewPostgresRepository creates a new postgres repository for sources
func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{
		db: d,
	}
}

func (r *postgresRepository) Create(ctx context.Context, s *Source) (*Source, error) {
	row, err := r.db.Queries.CreateSource(ctx, db.CreateSourceParams{
		Title:       s.Title,
		Subtitle:    r.toText(s.Subtitle),
		Type:        string(s.Type),
		Description: r.toText(s.Description),
		AuthorID:    r.toUUID(s.AuthorID),
		Publisher:   r.toText(s.Publisher),
		Isbn:        r.toText(s.ISBN),
		Doi:         r.toText(s.DOI),
		Url:         r.toText(s.URL),
		ExternalID:  r.toText(s.ExternalID),
		Tags:        r.toRawJSON(s.Tags),
		PublishedAt: r.toTimestamp(s.PublishedAt),
	})
	if err != nil {
		return nil, err
	}

	return r.mapToSource(row), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Source, error) {
	row, err := r.db.Queries.GetSource(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapToSource(row), nil
}

func (r *postgresRepository) List(ctx context.Context, limit, offset int) ([]*Source, error) {
	rows, err := r.db.Queries.ListSources(ctx, db.ListSourcesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	sources := make([]*Source, len(rows))
	for i, row := range rows {
		sources[i] = r.mapToSource(row)
	}
	return sources, nil
}

func (r *postgresRepository) ListByType(ctx context.Context, sourceType SourceType, limit, offset int) ([]*Source, error) {
	rows, err := r.db.Queries.ListSourcesByType(ctx, db.ListSourcesByTypeParams{
		Type:   string(sourceType),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	sources := make([]*Source, len(rows))
	for i, row := range rows {
		sources[i] = r.mapToSource(row)
	}
	return sources, nil
}

func (r *postgresRepository) Update(ctx context.Context, s *Source) (*Source, error) {
	row, err := r.db.Queries.UpdateSource(ctx, db.UpdateSourceParams{
		ID:          s.ID,
		Title:       r.toText(&s.Title),
		Subtitle:    r.toText(s.Subtitle),
		Type:        r.toText((*string)(&s.Type)),
		Description: r.toText(s.Description),
		AuthorID:    r.toUUID(s.AuthorID),
		Publisher:   r.toText(s.Publisher),
		Isbn:        r.toText(s.ISBN),
		Doi:         r.toText(s.DOI),
		Url:         r.toText(s.URL),
		ExternalID:  r.toText(s.ExternalID),
		Tags:        r.toRawJSON(s.Tags),
		PublishedAt: r.toTimestamp(s.PublishedAt),
	})
	if err != nil {
		return nil, err
	}

	return r.mapToSource(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Queries.DeleteSource(ctx, id)
}

func (r *postgresRepository) Search(ctx context.Context, query string, limit, offset int) ([]*Source, error) {
	rows, err := r.db.Queries.SearchSourcesByTitle(ctx, db.SearchSourcesByTitleParams{
		Column1: r.toText(&query),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	sources := make([]*Source, len(rows))
	for i, row := range rows {
		sources[i] = r.mapToSource(row)
	}
	return sources, nil
}

func (r *postgresRepository) Count(ctx context.Context) (int64, error) {
	return r.db.Queries.CountSources(ctx)
}

func (r *postgresRepository) mapToSource(row db.Source) *Source {
	var tags []string
	if row.Tags != nil {
		json.Unmarshal(row.Tags, &tags)
	}

	return &Source{
		ID:          row.ID,
		Title:       row.Title,
		Subtitle:    r.fromText(row.Subtitle),
		Type:        SourceType(row.Type),
		Description: r.fromText(row.Description),
		AuthorID:    r.fromUUID(row.AuthorID),
		Publisher:   r.fromText(row.Publisher),
		ISBN:        r.fromText(row.Isbn),
		DOI:         r.fromText(row.Doi),
		URL:         r.fromText(row.Url),
		ExternalID:  r.fromText(row.ExternalID),
		Tags:        tags,
		PublishedAt: r.fromTimestamp(row.PublishedAt),
		CreatedAt:   *r.fromTimestamp(row.CreatedAt),
		UpdatedAt:   *r.fromTimestamp(row.UpdatedAt),
	}
}

// Helper methods for pgtype conversions
func (r *postgresRepository) toText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func (r *postgresRepository) fromText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func (r *postgresRepository) toUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func (r *postgresRepository) fromUUID(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	res := uuid.UUID(u.Bytes)
	return &res
}

func (r *postgresRepository) toTimestamp(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func (r *postgresRepository) fromTimestamp(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func (r *postgresRepository) toRawJSON(data interface{}) []byte {
	if data == nil {
		return nil
	}
	b, _ := json.Marshal(data)
	return b
}
