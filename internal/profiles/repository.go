package profiles

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

type profileRow struct {
	ID            pgtype.UUID
	UserID        pgtype.UUID
	Username      pgtype.Text
	DisplayName   pgtype.Text
	Bio           pgtype.Text
	PublicProfile bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{db: d}
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	row, err := scanProfile(r.db.QueryRow(ctx, `
		SELECT p.id, p.user_id, u.username::text, p.display_name, p.bio, p.public_profile, p.created_at, p.updated_at
		FROM profiles p
		JOIN users u ON u.id = p.user_id
		WHERE p.user_id = $1
		LIMIT 1
	`, userID.String()))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func (r *postgresRepository) GetPublicByUsername(ctx context.Context, username string) (*Profile, error) {
	row, err := scanProfile(r.db.QueryRow(ctx, `
		SELECT p.id, p.user_id, u.username::text, p.display_name, p.bio, p.public_profile, p.created_at, p.updated_at
		FROM profiles p
		JOIN users u ON u.id = p.user_id
		WHERE u.username = $1 AND p.public_profile = true
		LIMIT 1
	`, username))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
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

	row, err := scanProfile(r.db.QueryRow(ctx, `
		WITH upserted AS (
			INSERT INTO profiles (id, user_id, display_name, bio, public_profile)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (user_id) DO UPDATE SET
				display_name = EXCLUDED.display_name,
				bio = EXCLUDED.bio,
				public_profile = EXCLUDED.public_profile,
				updated_at = NOW()
			RETURNING id, user_id, display_name, bio, public_profile, created_at, updated_at
		)
		SELECT p.id, p.user_id, u.username::text, p.display_name, p.bio, p.public_profile, p.created_at, p.updated_at
		FROM upserted p
		JOIN users u ON u.id = p.user_id
	`, id.String(), profile.UserID.String(), profile.DisplayName, profile.Bio, profile.PublicProfile))
	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func scanProfile(row pgx.Row) (profileRow, error) {
	var profile profileRow
	err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Username,
		&profile.DisplayName,
		&profile.Bio,
		&profile.PublicProfile,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	return profile, err
}

func mapRow(row profileRow) *Profile {
	return &Profile{
		ID:            uuid.UUID(row.ID.Bytes),
		UserID:        uuid.UUID(row.UserID.Bytes),
		Username:      row.Username.String,
		DisplayName:   stringPtr(row.DisplayName),
		Bio:           stringPtr(row.Bio),
		PublicProfile: row.PublicProfile,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func stringPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
