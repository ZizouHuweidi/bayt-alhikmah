package auth

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zizouhuweidi/maktaba/internal/db"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	GetUserByEmailOrUsername(ctx context.Context, login string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	CreateRefreshToken(ctx context.Context, token RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenHash []byte) (*RefreshToken, error)
	RotateRefreshToken(ctx context.Context, rotation RefreshTokenRotation) error
	RevokeRefreshTokenFamily(ctx context.Context, familyID uuid.UUID) error
}

type postgresRepository struct {
	db *db.DB
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{db: d}
}

func (r *postgresRepository) CreateUser(ctx context.Context, user User) (*User, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO users (id, email, username, password_hash)
		VALUES ($1, LOWER($2), LOWER($3), $4)
		RETURNING id, email, username, password_hash, created_at, updated_at
	`, user.ID.String(), user.Email, user.Username, user.PasswordHash)
	return scanUser(row)
}

func (r *postgresRepository) GetUserByEmailOrUsername(ctx context.Context, login string) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users WHERE email = LOWER($1) OR username = LOWER($1) LIMIT 1
	`, login)
	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

func (r *postgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users WHERE id = $1 LIMIT 1
	`, id.String())
	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

func (r *postgresRepository) CreateRefreshToken(ctx context.Context, token RefreshToken) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token_hash, family_id, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, token.ID.String(), token.UserID.String(), token.TokenHash, token.FamilyID.String(), token.ExpiresAt)
	return err
}

func (r *postgresRepository) GetRefreshToken(ctx context.Context, tokenHash []byte) (*RefreshToken, error) {
	var token RefreshToken
	var id pgtype.UUID
	var userID pgtype.UUID
	var familyID pgtype.UUID
	var revokedAt pgtype.Timestamptz
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, token_hash, family_id, expires_at, revoked_at, created_at
		FROM refresh_tokens WHERE token_hash = $1 LIMIT 1
	`, tokenHash).Scan(&id, &userID, &token.TokenHash, &familyID, &token.ExpiresAt, &revokedAt, &token.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token.ID = uuid.UUID(id.Bytes)
	token.UserID = uuid.UUID(userID.Bytes)
	token.FamilyID = uuid.UUID(familyID.Bytes)
	if revokedAt.Valid {
		t := revokedAt.Time
		token.RevokedAt = &t
	}
	return &token, nil
}

func (r *postgresRepository) RotateRefreshToken(ctx context.Context, rotation RefreshTokenRotation) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token_hash, family_id, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, rotation.NewToken.ID.String(), rotation.NewToken.UserID.String(), rotation.NewToken.TokenHash, rotation.NewToken.FamilyID.String(), rotation.NewToken.ExpiresAt)
	if err != nil {
		return err
	}

	result, err := tx.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = $2, replaced_by_token_id = $3
		WHERE id = $1 AND revoked_at IS NULL
	`, rotation.CurrentTokenID.String(), time.Now().UTC(), rotation.ReplacedByID.String())
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return ErrInvalidRefresh
	}

	return tx.Commit(ctx)
}

func (r *postgresRepository) RevokeRefreshTokenFamily(ctx context.Context, familyID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = COALESCE(revoked_at, $2) WHERE family_id = $1
	`, familyID.String(), time.Now().UTC())
	return err
}

func scanUser(row pgx.Row) (*User, error) {
	var user User
	var id pgtype.UUID
	if err := row.Scan(&id, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	user.ID = uuid.UUID(id.Bytes)
	return &user, nil
}
