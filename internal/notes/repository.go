package notes

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

type noteRow struct {
	ID          pgtype.UUID
	UserID      pgtype.UUID
	SourceID    pgtype.UUID
	Content     string
	ContentType string
	IsPublic    bool
	Annotations []byte
	Tags        []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{db: d}
}

func (r *postgresRepository) Create(ctx context.Context, n *Note) (*Note, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	annotations, err := json.Marshal(n.Annotations)
	if err != nil {
		return nil, err
	}
	tags, err := json.Marshal(n.Tags)
	if err != nil {
		return nil, err
	}

	row, err := scanNote(r.db.QueryRow(ctx, `
		INSERT INTO notes (id, user_id, source_id, content, content_type, is_public, annotations, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, source_id, content, content_type, is_public, annotations, tags, created_at, updated_at
	`, id.String(), n.UserID.String(), nullableUUIDString(n.SourceID), n.Content, string(n.ContentType), n.IsPublic, annotations, tags))
	if err != nil {
		return nil, err
	}

	return r.mapRow(row), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Note, error) {
	row, err := scanNote(r.db.QueryRow(ctx, `
		SELECT id, user_id, source_id, content, content_type, is_public, annotations, tags, created_at, updated_at
		FROM notes WHERE id = $1 LIMIT 1
	`, id.String()))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.mapRow(row), nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Note, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, source_id, content, content_type, is_public, annotations, tags, created_at, updated_at
		FROM notes WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *postgresRepository) ListBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Note, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, source_id, content, content_type, is_public, annotations, tags, created_at, updated_at
		FROM notes WHERE source_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, sourceID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *postgresRepository) ListPublic(ctx context.Context, limit, offset int) ([]*Note, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, source_id, content, content_type, is_public, annotations, tags, created_at, updated_at
		FROM notes WHERE is_public = true ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *postgresRepository) Update(ctx context.Context, n *Note) (*Note, error) {
	annotations, err := json.Marshal(n.Annotations)
	if err != nil {
		return nil, err
	}
	tags, err := json.Marshal(n.Tags)
	if err != nil {
		return nil, err
	}

	row, err := scanNote(r.db.QueryRow(ctx, `
		UPDATE notes
		SET content = $2, content_type = $3, is_public = $4, annotations = $5, tags = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, source_id, content, content_type, is_public, annotations, tags, created_at, updated_at
	`, n.ID.String(), n.Content, string(n.ContentType), n.IsPublic, annotations, tags))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.mapRow(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM notes WHERE id = $1`, id.String())
	return err
}

func (r *postgresRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM notes WHERE user_id = $1`, userID.String()).Scan(&count)
	return count, err
}

func (r *postgresRepository) scanRows(rows pgx.Rows) ([]*Note, error) {
	var notes []*Note
	for rows.Next() {
		row, err := scanNote(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, r.mapRow(row))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func scanNote(row pgx.Row) (noteRow, error) {
	var n noteRow
	err := row.Scan(
		&n.ID,
		&n.UserID,
		&n.SourceID,
		&n.Content,
		&n.ContentType,
		&n.IsPublic,
		&n.Annotations,
		&n.Tags,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	return n, err
}

func (r *postgresRepository) mapRow(row noteRow) *Note {
	var annotations []string
	if len(row.Annotations) > 0 {
		_ = json.Unmarshal(row.Annotations, &annotations)
	}
	var tags []string
	if len(row.Tags) > 0 {
		_ = json.Unmarshal(row.Tags, &tags)
	}

	return &Note{
		ID:          uuid.UUID(row.ID.Bytes),
		UserID:      uuid.UUID(row.UserID.Bytes),
		SourceID:    uuidPtr(row.SourceID),
		Content:     row.Content,
		ContentType: ContentType(row.ContentType),
		IsPublic:    row.IsPublic,
		Annotations: annotations,
		Tags:        tags,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
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
