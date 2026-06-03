package reviews

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

func (r *postgresRepository) Create(ctx context.Context, review *Review) (*Review, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	row, err := r.queries.CreateReview(ctx, dbgen.CreateReviewParams{
		ID:       db.PGUUID(id),
		UserID:   db.PGUUID(review.UserID),
		SourceID: db.PGUUID(review.SourceID),
		Rating:   pgtype.Int4{Int32: int32(review.Rating), Valid: true},
		Content:  db.PGText(review.Content),
		IsPublic: db.PGBool(review.IsPublic),
	})
	if err != nil {
		return nil, mapCreateError(err)
	}
	return mapReview(row), nil
}

func mapCreateError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505":
		return ErrReviewExists
	case "23503":
		if pgErr.ConstraintName == "reviews_source_id_fkey" {
			return ErrSourceNotFound
		}
		return ErrReviewConflict
	default:
		return err
	}
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Review, error) {
	row, err := r.queries.GetReviewByID(ctx, db.PGUUID(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapReview(row), nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Review, error) {
	rows, err := r.queries.ListReviewsByUser(ctx, dbgen.ListReviewsByUserParams{UserID: db.PGUUID(userID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapReviews(rows), nil
}

func (r *postgresRepository) ListBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Review, error) {
	rows, err := r.queries.ListReviewsBySource(ctx, dbgen.ListReviewsBySourceParams{SourceID: db.PGUUID(sourceID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapReviews(rows), nil
}

func (r *postgresRepository) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Review, error) {
	rows, err := r.queries.ListPublicReviewsByUser(ctx, dbgen.ListPublicReviewsByUserParams{UserID: db.PGUUID(userID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapReviews(rows), nil
}

func (r *postgresRepository) ListPublicBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Review, error) {
	rows, err := r.queries.ListPublicReviewsBySource(ctx, dbgen.ListPublicReviewsBySourceParams{SourceID: db.PGUUID(sourceID), Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return mapReviews(rows), nil
}

func (r *postgresRepository) Update(ctx context.Context, review *Review) (*Review, error) {
	row, err := r.queries.UpdateReview(ctx, dbgen.UpdateReviewParams{ID: db.PGUUID(review.ID), Rating: pgtype.Int4{Int32: int32(review.Rating), Valid: true}, Content: db.PGText(review.Content), IsPublic: db.PGBool(review.IsPublic)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapReview(row), nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteReview(ctx, db.PGUUID(id))
}

func mapReviews(rows []dbgen.Review) []*Review {
	reviews := make([]*Review, 0, len(rows))
	for _, row := range rows {
		reviews = append(reviews, mapReview(row))
	}
	return reviews
}

func mapReview(row dbgen.Review) *Review {
	return &Review{
		ID:        db.UUID(row.ID),
		UserID:    db.UUID(row.UserID),
		SourceID:  db.UUID(row.SourceID),
		Rating:    int(row.Rating.Int32),
		Content:   db.StringPtr(row.Content),
		IsPublic:  db.Bool(row.IsPublic),
		CreatedAt: db.Time(row.CreatedAt),
		UpdatedAt: db.Time(row.UpdatedAt),
	}
}
