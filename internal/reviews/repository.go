package reviews

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zizouhuweidi/maktaba/internal/db"
)

type postgresRepository struct {
	db *db.DB
}

type reviewRow struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	SourceID  pgtype.UUID
	Rating    int
	Content   pgtype.Text
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{db: d}
}

func (r *postgresRepository) Create(ctx context.Context, review *Review) (*Review, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	row, err := scanReview(r.db.QueryRow(ctx, `
		INSERT INTO reviews (id, user_id, source_id, rating, content, is_public)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, source_id, rating, content, is_public, created_at, updated_at
	`, id.String(), review.UserID.String(), review.SourceID.String(), review.Rating, review.Content, review.IsPublic))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrReviewExists
		}
		return nil, err
	}

	return mapRow(row), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Review, error) {
	row, err := scanReview(r.db.QueryRow(ctx, `
		SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
		FROM reviews WHERE id = $1 LIMIT 1
	`, id.String()))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Review, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
		FROM reviews WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (r *postgresRepository) ListBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Review, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
		FROM reviews WHERE source_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, sourceID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (r *postgresRepository) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Review, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
		FROM reviews WHERE user_id = $1 AND is_public = true ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (r *postgresRepository) ListPublicBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Review, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, source_id, rating, content, is_public, created_at, updated_at
		FROM reviews WHERE source_id = $1 AND is_public = true ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, sourceID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (r *postgresRepository) Update(ctx context.Context, review *Review) (*Review, error) {
	row, err := scanReview(r.db.QueryRow(ctx, `
		UPDATE reviews
		SET rating = $2, content = $3, is_public = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, source_id, rating, content, is_public, created_at, updated_at
	`, review.ID.String(), review.Rating, review.Content, review.IsPublic))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM reviews WHERE id = $1`, id.String())
	return err
}

func scanRows(rows pgx.Rows) ([]*Review, error) {
	var reviews []*Review
	for rows.Next() {
		row, err := scanReview(rows)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, mapRow(row))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}

func scanReview(row pgx.Row) (reviewRow, error) {
	var review reviewRow
	err := row.Scan(
		&review.ID,
		&review.UserID,
		&review.SourceID,
		&review.Rating,
		&review.Content,
		&review.IsPublic,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	return review, err
}

func mapRow(row reviewRow) *Review {
	return &Review{
		ID:        uuid.UUID(row.ID.Bytes),
		UserID:    uuid.UUID(row.UserID.Bytes),
		SourceID:  uuid.UUID(row.SourceID.Bytes),
		Rating:    row.Rating,
		Content:   stringPtr(row.Content),
		IsPublic:  row.IsPublic,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func stringPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
