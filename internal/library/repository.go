package library

import (
	"context"
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

type itemRow struct {
	ID            pgtype.UUID
	UserID        pgtype.UUID
	SourceID      pgtype.UUID
	Status        string
	ProgressValue pgtype.Int4
	ProgressUnit  pgtype.Text
	Visibility    string
	StartedAt     pgtype.Timestamptz
	CompletedAt   pgtype.Timestamptz
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{db: d}
}

func (r *postgresRepository) Create(ctx context.Context, item *Item) (*Item, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	row, err := scanItem(r.db.QueryRow(ctx, `
		INSERT INTO user_library_items (id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at
	`, id.String(), item.UserID.String(), item.SourceID.String(), string(item.Status), item.ProgressValue, progressUnitString(item.ProgressUnit), string(item.Visibility), item.StartedAt, item.CompletedAt))
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Item, error) {
	row, err := scanItem(r.db.QueryRow(ctx, `
		SELECT id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at
		FROM user_library_items WHERE id = $1 LIMIT 1
	`, id.String()))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at
		FROM user_library_items WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanItems(rows)
}

func (r *postgresRepository) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at
		FROM user_library_items WHERE user_id = $1 AND visibility = 'public' ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanItems(rows)
}

func (r *postgresRepository) Update(ctx context.Context, item *Item) (*Item, error) {
	row, err := scanItem(r.db.QueryRow(ctx, `
		UPDATE user_library_items
		SET status = $2, progress_value = $3, progress_unit = $4, visibility = $5, started_at = $6, completed_at = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, source_id, status, progress_value, progress_unit, visibility, started_at, completed_at, created_at, updated_at
	`, item.ID.String(), string(item.Status), item.ProgressValue, progressUnitString(item.ProgressUnit), string(item.Visibility), item.StartedAt, item.CompletedAt))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM user_library_items WHERE id = $1`, id.String())
	return err
}

func scanItems(rows pgx.Rows) ([]*Item, error) {
	var items []*Item
	for rows.Next() {
		row, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, mapRow(row))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func scanItem(row pgx.Row) (itemRow, error) {
	var item itemRow
	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.SourceID,
		&item.Status,
		&item.ProgressValue,
		&item.ProgressUnit,
		&item.Visibility,
		&item.StartedAt,
		&item.CompletedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}

func mapRow(row itemRow) *Item {
	return &Item{
		ID:            uuid.UUID(row.ID.Bytes),
		UserID:        uuid.UUID(row.UserID.Bytes),
		SourceID:      uuid.UUID(row.SourceID.Bytes),
		Status:        Status(row.Status),
		ProgressValue: intPtr(row.ProgressValue),
		ProgressUnit:  progressUnitPtr(row.ProgressUnit),
		Visibility:    Visibility(row.Visibility),
		StartedAt:     timePtr(row.StartedAt),
		CompletedAt:   timePtr(row.CompletedAt),
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func intPtr(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	result := int(value.Int32)
	return &result
}

func progressUnitPtr(value pgtype.Text) *ProgressUnit {
	if !value.Valid {
		return nil
	}
	unit := ProgressUnit(value.String)
	return &unit
}

func progressUnitString(value *ProgressUnit) *string {
	if value == nil {
		return nil
	}
	result := string(*value)
	return &result
}

func timePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}
