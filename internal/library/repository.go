package library

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zizouhuweidi/maktaba/internal/db"
	"github.com/zizouhuweidi/maktaba/internal/db/dbgen"
)

type postgresRepository struct {
	queries *dbgen.Queries
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{queries: dbgen.New(d.Pool)}
}

func (r *postgresRepository) Create(ctx context.Context, item *Item) (*Item, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	row, err := r.queries.CreateLibraryItem(ctx, dbgen.CreateLibraryItemParams{
		ID:            db.PGUUID(id),
		UserID:        db.PGUUID(item.UserID),
		SourceID:      db.PGUUID(item.SourceID),
		Status:        string(item.Status),
		ProgressValue: db.PGInt4Ptr(item.ProgressValue),
		ProgressUnit:  pgProgressUnit(item.ProgressUnit),
		Visibility:    string(item.Visibility),
		StartedAt:     db.PGTimestamptzPtr(item.StartedAt),
		CompletedAt:   db.PGTimestamptzPtr(item.CompletedAt),
	})
	if err != nil {
		return nil, mapCreateError(err)
	}
	return mapItem(row), nil
}

func mapCreateError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505":
		return ErrItemExists
	case "23503":
		if pgErr.ConstraintName == "user_library_items_source_id_fkey" {
			return ErrSourceNotFound
		}
		return ErrLibraryConflict
	default:
		return err
	}
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Item, error) {
	row, err := r.queries.GetLibraryItemByID(ctx, db.PGUUID(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapItem(row), nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error) {
	rows, err := r.queries.ListLibraryItemsByUser(ctx, dbgen.ListLibraryItemsByUserParams{UserID: db.PGUUID(userID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapItems(rows), nil
}

func (r *postgresRepository) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error) {
	rows, err := r.queries.ListPublicLibraryItemsByUser(ctx, dbgen.ListPublicLibraryItemsByUserParams{UserID: db.PGUUID(userID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapItems(rows), nil
}

func (r *postgresRepository) ListPublicByUsername(ctx context.Context, username string, limit, offset int) ([]*Item, error) {
	rows, err := r.queries.ListPublicLibraryItemsByUsername(ctx, dbgen.ListPublicLibraryItemsByUsernameParams{Username: username, Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapItems(rows), nil
}

func (r *postgresRepository) ListByUserWithSources(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ItemWithSource, error) {
	rows, err := r.queries.ListLibraryItemsByUserWithSources(ctx, dbgen.ListLibraryItemsByUserWithSourcesParams{UserID: db.PGUUID(userID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	items := make([]*ItemWithSource, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapItemWithSource(row))
	}
	return items, nil
}

func (r *postgresRepository) ListPublicByUsernameWithSources(ctx context.Context, username string, limit, offset int) ([]*ItemWithSource, error) {
	rows, err := r.queries.ListPublicLibraryItemsByUsernameWithSources(ctx, dbgen.ListPublicLibraryItemsByUsernameWithSourcesParams{Username: username, Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	items := make([]*ItemWithSource, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPublicItemWithSource(row))
	}
	return items, nil
}

func (r *postgresRepository) Update(ctx context.Context, item *Item) (*Item, error) {
	row, err := r.queries.UpdateLibraryItem(ctx, dbgen.UpdateLibraryItemParams{
		ID:            db.PGUUID(item.ID),
		Status:        string(item.Status),
		ProgressValue: db.PGInt4Ptr(item.ProgressValue),
		ProgressUnit:  pgProgressUnit(item.ProgressUnit),
		Visibility:    string(item.Visibility),
		StartedAt:     db.PGTimestamptzPtr(item.StartedAt),
		CompletedAt:   db.PGTimestamptzPtr(item.CompletedAt),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapItem(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteLibraryItem(ctx, db.PGUUID(id))
}

func mapItems(rows []dbgen.UserLibraryItem) []*Item {
	items := make([]*Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapItem(row))
	}
	return items
}

func mapItem(row dbgen.UserLibraryItem) *Item {
	return &Item{
		ID:            db.UUID(row.ID),
		UserID:        db.UUID(row.UserID),
		SourceID:      db.UUID(row.SourceID),
		Status:        Status(row.Status),
		ProgressValue: db.IntPtr(row.ProgressValue),
		ProgressUnit:  progressUnitPtr(row.ProgressUnit),
		Visibility:    Visibility(row.Visibility),
		StartedAt:     db.TimePtr(row.StartedAt),
		CompletedAt:   db.TimePtr(row.CompletedAt),
		CreatedAt:     db.Time(row.CreatedAt),
		UpdatedAt:     db.Time(row.UpdatedAt),
	}
}

func mapItemWithSource(row dbgen.ListLibraryItemsByUserWithSourcesRow) *ItemWithSource {
	item := mapJoinedItem(row.ID, row.UserID, row.SourceID, row.Status, row.ProgressValue, row.ProgressUnit, row.Visibility, row.StartedAt, row.CompletedAt, row.CreatedAt, row.UpdatedAt)
	return &ItemWithSource{Item: item, Source: mapSourceSummary(item.SourceID, row.Title, row.Subtitle, row.Type, row.Publisher, row.Isbn)}
}

func mapPublicItemWithSource(row dbgen.ListPublicLibraryItemsByUsernameWithSourcesRow) *ItemWithSource {
	item := mapJoinedItem(row.ID, row.UserID, row.SourceID, row.Status, row.ProgressValue, row.ProgressUnit, row.Visibility, row.StartedAt, row.CompletedAt, row.CreatedAt, row.UpdatedAt)
	return &ItemWithSource{Item: item, Source: mapSourceSummary(item.SourceID, row.Title, row.Subtitle, row.Type, row.Publisher, row.Isbn)}
}

func mapJoinedItem(id, userID, sourceID pgtype.UUID, status string, progressValue pgtype.Int4, progressUnit pgtype.Text, visibility string, startedAt, completedAt, createdAt, updatedAt pgtype.Timestamptz) *Item {
	return &Item{
		ID:            db.UUID(id),
		UserID:        db.UUID(userID),
		SourceID:      db.UUID(sourceID),
		Status:        Status(status),
		ProgressValue: db.IntPtr(progressValue),
		ProgressUnit:  progressUnitPtr(progressUnit),
		Visibility:    Visibility(visibility),
		StartedAt:     db.TimePtr(startedAt),
		CompletedAt:   db.TimePtr(completedAt),
		CreatedAt:     db.Time(createdAt),
		UpdatedAt:     db.Time(updatedAt),
	}
}

func mapSourceSummary(id uuid.UUID, title string, subtitle pgtype.Text, sourceType string, publisher pgtype.Text, isbn pgtype.Text) *SourceSummary {
	return &SourceSummary{
		ID:        id,
		Title:     title,
		Subtitle:  db.StringPtr(subtitle),
		Type:      sourceType,
		Publisher: db.StringPtr(publisher),
		ISBN:      db.StringPtr(isbn),
	}
}

func pgProgressUnit(value *ProgressUnit) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: string(*value), Valid: true}
}

func progressUnitPtr(value pgtype.Text) *ProgressUnit {
	if !value.Valid {
		return nil
	}
	unit := ProgressUnit(value.String)
	return &unit
}
