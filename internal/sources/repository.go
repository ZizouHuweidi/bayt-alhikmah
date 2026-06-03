package sources

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zizouhuweidi/maktaba/internal/db"
	"github.com/zizouhuweidi/maktaba/internal/db/dbgen"
)

type postgresRepository struct {
	db      *db.DB
	queries *dbgen.Queries
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{db: d, queries: dbgen.New(d.Pool)}
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

	row, err := r.queries.CreateSource(ctx, dbgen.CreateSourceParams{
		ID:          db.PGUUID(id),
		Title:       s.Title,
		Subtitle:    db.PGText(s.Subtitle),
		Type:        string(s.Type),
		Description: db.PGText(s.Description),
		Publisher:   db.PGText(s.Publisher),
		Isbn:        db.PGText(s.ISBN),
		Doi:         db.PGText(s.DOI),
		Url:         db.PGText(s.URL),
		ExternalID:  db.PGText(s.ExternalID),
		Tags:        tags,
		PublishedAt: db.PGTimestamptzPtr(s.PublishedAt),
	})
	if err != nil {
		return nil, err
	}
	return mapSource(row), nil
}

func (r *postgresRepository) CreateBook(ctx context.Context, params CreateBookParams) (*Book, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)

	sourceID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	tags, err := json.Marshal(params.Tags)
	if err != nil {
		return nil, err
	}

	sourceRow, err := qtx.InsertBookSource(ctx, dbgen.InsertBookSourceParams{
		ID:          db.PGUUID(sourceID),
		Title:       params.Title,
		Subtitle:    db.PGText(params.Subtitle),
		Description: db.PGText(params.Description),
		Publisher:   db.PGText(params.Publisher),
		Column6:     db.PGText(params.ISBN13),
		Column7:     db.PGText(params.ISBN10),
		Url:         db.PGText(params.URL),
		ExternalID:  db.PGText(params.ExternalID),
		Tags:        tags,
		PublishedAt: db.PGTimestamptzPtr(params.PublishedAt),
	})
	if err != nil {
		return nil, err
	}

	metadataRow, err := qtx.InsertBookMetadata(ctx, dbgen.InsertBookMetadataParams{
		SourceID:  db.PGUUID(sourceID),
		Isbn10:    db.PGText(params.ISBN10),
		Isbn13:    db.PGText(params.ISBN13),
		Publisher: db.PGText(params.Publisher),
		PageCount: db.PGInt4Ptr(params.PageCount),
		Language:  db.PGText(params.Language),
		CoverUrl:  db.PGText(params.CoverURL),
	})
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
		existingID, err := qtx.UpsertContributor(ctx, dbgen.UpsertContributorParams{ID: db.PGUUID(contributorID), Name: contributor.Name})
		if err != nil {
			return nil, err
		}
		row, err := qtx.InsertSourceContributor(ctx, dbgen.InsertSourceContributorParams{
			SourceID:        db.PGUUID(sourceID),
			ContributorID:   existingID,
			Role:            role,
			Position:        int32(position),
			ContributorName: contributor.Name,
		})
		if err != nil {
			return nil, err
		}
		contributors = append(contributors, mapInsertedContributor(row))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &Book{Source: mapSource(sourceRow), Metadata: mapBookMetadata(metadataRow), Contributors: contributors}, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Source, error) {
	row, err := r.queries.GetSourceByID(ctx, db.PGUUID(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapSource(row), nil
}

func (r *postgresRepository) GetBookByID(ctx context.Context, id uuid.UUID) (*Book, error) {
	sourceRow, err := r.queries.GetBookSourceByID(ctx, db.PGUUID(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	metadataRow, err := r.queries.GetBookMetadata(ctx, db.PGUUID(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return &Book{Source: mapSource(sourceRow)}, nil
	}
	if err != nil {
		return nil, err
	}
	contributors, err := r.listContributors(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Book{Source: mapSource(sourceRow), Metadata: mapBookMetadata(metadataRow), Contributors: contributors}, nil
}

func (r *postgresRepository) List(ctx context.Context, limit, offset int) ([]*Source, error) {
	rows, err := r.queries.ListSources(ctx, dbgen.ListSourcesParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapSources(rows), nil
}

func (r *postgresRepository) ListByType(ctx context.Context, sourceType SourceType, limit, offset int) ([]*Source, error) {
	rows, err := r.queries.ListSourcesByType(ctx, dbgen.ListSourcesByTypeParams{Type: string(sourceType), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapSources(rows), nil
}

func (r *postgresRepository) Update(ctx context.Context, s *Source) (*Source, error) {
	tags, err := json.Marshal(s.Tags)
	if err != nil {
		return nil, err
	}
	row, err := r.queries.UpdateSource(ctx, dbgen.UpdateSourceParams{
		ID:          db.PGUUID(s.ID),
		Title:       s.Title,
		Subtitle:    db.PGText(s.Subtitle),
		Type:        string(s.Type),
		Description: db.PGText(s.Description),
		Publisher:   db.PGText(s.Publisher),
		Isbn:        db.PGText(s.ISBN),
		Doi:         db.PGText(s.DOI),
		Url:         db.PGText(s.URL),
		ExternalID:  db.PGText(s.ExternalID),
		Tags:        tags,
		PublishedAt: db.PGTimestamptzPtr(s.PublishedAt),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapSource(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSource(ctx, db.PGUUID(id))
}

func (r *postgresRepository) Search(ctx context.Context, query string, limit, offset int) ([]*Source, error) {
	rows, err := r.queries.SearchSources(ctx, dbgen.SearchSourcesParams{Column1: pgtype.Text{String: query, Valid: true}, Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapSources(rows), nil
}

func (r *postgresRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountSources(ctx)
}

func (r *postgresRepository) listContributors(ctx context.Context, sourceID uuid.UUID) ([]*Contributor, error) {
	rows, err := r.queries.ListContributorsBySource(ctx, db.PGUUID(sourceID))
	if err != nil {
		return nil, err
	}
	contributors := make([]*Contributor, 0, len(rows))
	for _, row := range rows {
		contributors = append(contributors, mapContributor(row))
	}
	return contributors, nil
}

func mapSources(rows []dbgen.Source) []*Source {
	sources := make([]*Source, 0, len(rows))
	for _, row := range rows {
		sources = append(sources, mapSource(row))
	}
	return sources
}

func mapSource(row dbgen.Source) *Source {
	var tags []string
	if len(row.Tags) > 0 {
		_ = json.Unmarshal(row.Tags, &tags)
	}
	return &Source{
		ID:          db.UUID(row.ID),
		Title:       row.Title,
		Subtitle:    db.StringPtr(row.Subtitle),
		Type:        SourceType(row.Type),
		Description: db.StringPtr(row.Description),
		Publisher:   db.StringPtr(row.Publisher),
		ISBN:        db.StringPtr(row.Isbn),
		DOI:         db.StringPtr(row.Doi),
		URL:         db.StringPtr(row.Url),
		ExternalID:  db.StringPtr(row.ExternalID),
		Tags:        tags,
		PublishedAt: db.TimePtr(row.PublishedAt),
		CreatedAt:   db.Time(row.CreatedAt),
		UpdatedAt:   db.Time(row.UpdatedAt),
	}
}

func mapBookMetadata(row dbgen.BookMetadatum) *BookMetadata {
	return &BookMetadata{
		SourceID:  db.UUID(row.SourceID),
		ISBN10:    db.StringPtr(row.Isbn10),
		ISBN13:    db.StringPtr(row.Isbn13),
		Publisher: db.StringPtr(row.Publisher),
		PageCount: db.IntPtr(row.PageCount),
		Language:  db.StringPtr(row.Language),
		CoverURL:  db.StringPtr(row.CoverUrl),
		CreatedAt: db.Time(row.CreatedAt),
		UpdatedAt: db.Time(row.UpdatedAt),
	}
}

func mapContributor(row dbgen.ListContributorsBySourceRow) *Contributor {
	return &Contributor{
		ID:        db.UUID(row.ID),
		Name:      row.Name,
		Role:      row.Role,
		Position:  int(row.Position),
		CreatedAt: db.Time(row.CreatedAt),
		UpdatedAt: db.Time(row.UpdatedAt),
	}
}

func mapInsertedContributor(row dbgen.InsertSourceContributorRow) *Contributor {
	return &Contributor{
		ID:        db.UUID(row.ContributorID),
		Name:      row.Name,
		Role:      row.Role,
		Position:  int(row.Position),
		CreatedAt: db.Time(row.CreatedAt),
		UpdatedAt: db.Time(row.UpdatedAt),
	}
}
