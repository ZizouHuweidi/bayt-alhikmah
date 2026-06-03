package collections

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/zizouhuweidi/maktaba/internal/db"
	"github.com/zizouhuweidi/maktaba/internal/db/dbgen"
)

type postgresRepository struct {
	queries *dbgen.Queries
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{queries: dbgen.New(d.Pool)}
}

func (r *postgresRepository) Create(ctx context.Context, collection *Collection) (*Collection, error) {
	if err := r.validateSourceIDs(ctx, collection.SourceIDs); err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	sourceIDs, err := json.Marshal(collection.SourceIDs)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.CreateCollection(ctx, dbgen.CreateCollectionParams{
		ID:          db.PGUUID(id),
		UserID:      db.PGUUID(collection.UserID),
		Name:        collection.Name,
		Description: db.PGText(collection.Description),
		IsPublic:    db.PGBool(collection.IsPublic),
		SourceIds:   sourceIDs,
	})
	if err != nil {
		return nil, err
	}
	return mapCollection(row), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Collection, error) {
	row, err := r.queries.GetCollectionByID(ctx, db.PGUUID(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapCollection(row), nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Collection, error) {
	rows, err := r.queries.ListCollectionsByUser(ctx, dbgen.ListCollectionsByUserParams{UserID: db.PGUUID(userID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapCollections(rows), nil
}

func (r *postgresRepository) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Collection, error) {
	rows, err := r.queries.ListPublicCollectionsByUser(ctx, dbgen.ListPublicCollectionsByUserParams{UserID: db.PGUUID(userID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapCollections(rows), nil
}

func (r *postgresRepository) Update(ctx context.Context, collection *Collection) (*Collection, error) {
	if err := r.validateSourceIDs(ctx, collection.SourceIDs); err != nil {
		return nil, err
	}

	sourceIDs, err := json.Marshal(collection.SourceIDs)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.UpdateCollection(ctx, dbgen.UpdateCollectionParams{
		ID:          db.PGUUID(collection.ID),
		Name:        collection.Name,
		Description: db.PGText(collection.Description),
		IsPublic:    db.PGBool(collection.IsPublic),
		SourceIds:   sourceIDs,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapCollection(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteCollection(ctx, db.PGUUID(id))
}

func (r *postgresRepository) validateSourceIDs(ctx context.Context, sourceIDs []uuid.UUID) error {
	for _, sourceID := range sourceIDs {
		exists, err := r.queries.SourceExists(ctx, db.PGUUID(sourceID))
		if err != nil {
			return err
		}
		if !exists {
			return ErrSourceNotFound
		}
	}
	return nil
}

func mapCollections(rows []dbgen.Collection) []*Collection {
	collections := make([]*Collection, 0, len(rows))
	for _, row := range rows {
		collections = append(collections, mapCollection(row))
	}
	return collections
}

func mapCollection(row dbgen.Collection) *Collection {
	var sourceIDs []uuid.UUID
	if len(row.SourceIds) > 0 {
		_ = json.Unmarshal(row.SourceIds, &sourceIDs)
	}

	return &Collection{
		ID:          db.UUID(row.ID),
		UserID:      db.UUID(row.UserID),
		Name:        row.Name,
		Description: db.StringPtr(row.Description),
		IsPublic:    db.Bool(row.IsPublic),
		SourceIDs:   sourceIDs,
		CreatedAt:   db.Time(row.CreatedAt),
		UpdatedAt:   db.Time(row.UpdatedAt),
	}
}
