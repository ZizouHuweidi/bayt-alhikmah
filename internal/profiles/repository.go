package profiles

import (
	"context"
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

func (r *postgresRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	row, err := r.queries.GetProfileByUserID(ctx, db.PGUUID(userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	profile := mapProfile(row)
	return profile, nil
}

func (r *postgresRepository) GetPublicByUsername(ctx context.Context, username string) (*Profile, error) {
	row, err := r.queries.GetPublicProfileByUsername(ctx, username)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &Profile{
		ID:            db.UUID(row.ID),
		UserID:        db.UUID(row.UserID),
		Username:      row.Username,
		DisplayName:   db.StringPtr(row.DisplayName),
		Bio:           db.StringPtr(row.Bio),
		PublicProfile: db.Bool(row.PublicProfile),
		CreatedAt:     db.Time(row.CreatedAt),
		UpdatedAt:     db.Time(row.UpdatedAt),
	}, nil
}

func (r *postgresRepository) Upsert(ctx context.Context, profile *Profile) (*Profile, error) {
	id := profile.ID
	if id == uuid.Nil {
		var err error
		id, err = uuid.NewV7()
		if err != nil {
			return nil, err
		}
	}

	row, err := r.queries.UpsertProfile(ctx, dbgen.UpsertProfileParams{
		ID:            db.PGUUID(id),
		UserID:        db.PGUUID(profile.UserID),
		DisplayName:   db.PGText(profile.DisplayName),
		Bio:           db.PGText(profile.Bio),
		PublicProfile: db.PGBool(profile.PublicProfile),
	})
	if err != nil {
		return nil, err
	}
	return &Profile{
		ID:            db.UUID(row.ID),
		UserID:        db.UUID(row.UserID),
		Username:      row.Username,
		DisplayName:   db.StringPtr(row.DisplayName),
		Bio:           db.StringPtr(row.Bio),
		PublicProfile: db.Bool(row.PublicProfile),
		CreatedAt:     db.Time(row.CreatedAt),
		UpdatedAt:     db.Time(row.UpdatedAt),
	}, nil
}

func mapProfile(row dbgen.Profile) *Profile {
	return &Profile{
		ID:            db.UUID(row.ID),
		UserID:        db.UUID(row.UserID),
		DisplayName:   db.StringPtr(row.DisplayName),
		Bio:           db.StringPtr(row.Bio),
		PublicProfile: db.Bool(row.PublicProfile),
		CreatedAt:     db.Time(row.CreatedAt),
		UpdatedAt:     db.Time(row.UpdatedAt),
	}
}
