package collections

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

type collectionRow struct {
	ID          pgtype.UUID
	UserID      pgtype.UUID
	Name        string
	Description pgtype.Text
	IsPublic    bool
	SourceIDs   []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{db: d}
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

	row, err := scanCollection(r.db.QueryRow(ctx, `
		INSERT INTO collections (id, user_id, name, description, is_public, source_ids)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, name, description, is_public, source_ids, created_at, updated_at
	`, id.String(), collection.UserID.String(), collection.Name, collection.Description, collection.IsPublic, sourceIDs))
	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Collection, error) {
	row, err := scanCollection(r.db.QueryRow(ctx, `
		SELECT id, user_id, name, description, is_public, source_ids, created_at, updated_at
		FROM collections WHERE id = $1 LIMIT 1
	`, id.String()))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Collection, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, name, description, is_public, source_ids, created_at, updated_at
		FROM collections WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (r *postgresRepository) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Collection, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, name, description, is_public, source_ids, created_at, updated_at
		FROM collections WHERE user_id = $1 AND is_public = true ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (r *postgresRepository) Update(ctx context.Context, collection *Collection) (*Collection, error) {
	if err := r.validateSourceIDs(ctx, collection.SourceIDs); err != nil {
		return nil, err
	}

	sourceIDs, err := json.Marshal(collection.SourceIDs)
	if err != nil {
		return nil, err
	}

	row, err := scanCollection(r.db.QueryRow(ctx, `
		UPDATE collections
		SET name = $2, description = $3, is_public = $4, source_ids = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, name, description, is_public, source_ids, created_at, updated_at
	`, collection.ID.String(), collection.Name, collection.Description, collection.IsPublic, sourceIDs))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM collections WHERE id = $1`, id.String())
	return err
}

func (r *postgresRepository) validateSourceIDs(ctx context.Context, sourceIDs []uuid.UUID) error {
	for _, sourceID := range sourceIDs {
		var exists bool
		err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM sources WHERE id = $1)`, sourceID.String()).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return ErrSourceNotFound
		}
	}
	return nil
}

func scanRows(rows pgx.Rows) ([]*Collection, error) {
	var collections []*Collection
	for rows.Next() {
		row, err := scanCollection(rows)
		if err != nil {
			return nil, err
		}
		collections = append(collections, mapRow(row))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return collections, nil
}

func scanCollection(row pgx.Row) (collectionRow, error) {
	var collection collectionRow
	err := row.Scan(
		&collection.ID,
		&collection.UserID,
		&collection.Name,
		&collection.Description,
		&collection.IsPublic,
		&collection.SourceIDs,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)
	return collection, err
}

func mapRow(row collectionRow) *Collection {
	var sourceIDs []uuid.UUID
	if len(row.SourceIDs) > 0 {
		_ = json.Unmarshal(row.SourceIDs, &sourceIDs)
	}

	return &Collection{
		ID:          uuid.UUID(row.ID.Bytes),
		UserID:      uuid.UUID(row.UserID.Bytes),
		Name:        row.Name,
		Description: stringPtr(row.Description),
		IsPublic:    row.IsPublic,
		SourceIDs:   sourceIDs,
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
