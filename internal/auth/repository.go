package auth

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/zizouhuweidi/maktaba/internal/db"
	"github.com/zizouhuweidi/maktaba/internal/db/dbgen"
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
	db      *db.DB
	queries *dbgen.Queries
}

func NewPostgresRepository(d *db.DB) Repository {
	return &postgresRepository{db: d, queries: dbgen.New(d.Pool)}
}

func (r *postgresRepository) CreateUser(ctx context.Context, user User) (*User, error) {
	row, err := r.queries.CreateUser(ctx, dbgen.CreateUserParams{
		ID:           db.PGUUID(user.ID),
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return nil, err
	}
	return mapCreateUserRow(row), nil
}

func (r *postgresRepository) GetUserByEmailOrUsername(ctx context.Context, login string) (*User, error) {
	row, err := r.queries.GetUserByEmailOrUsername(ctx, login)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapGetUserByEmailOrUsernameRow(row), nil
}

func (r *postgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row, err := r.queries.GetUserByID(ctx, db.PGUUID(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapGetUserByIDRow(row), nil
}

func (r *postgresRepository) CreateRefreshToken(ctx context.Context, token RefreshToken) error {
	return r.queries.CreateRefreshToken(ctx, dbgen.CreateRefreshTokenParams{
		ID:        db.PGUUID(token.ID),
		UserID:    db.PGUUID(token.UserID),
		TokenHash: token.TokenHash,
		FamilyID:  db.PGUUID(token.FamilyID),
		ExpiresAt: db.PGTimestamptz(token.ExpiresAt),
	})
}

func (r *postgresRepository) GetRefreshToken(ctx context.Context, tokenHash []byte) (*RefreshToken, error) {
	row, err := r.queries.GetRefreshToken(ctx, tokenHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapRefreshTokenRow(row), nil
}

func (r *postgresRepository) RotateRefreshToken(ctx context.Context, rotation RefreshTokenRotation) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)
	if err := qtx.InsertRotatedRefreshToken(ctx, dbgen.InsertRotatedRefreshTokenParams{
		ID:        db.PGUUID(rotation.NewToken.ID),
		UserID:    db.PGUUID(rotation.NewToken.UserID),
		TokenHash: rotation.NewToken.TokenHash,
		FamilyID:  db.PGUUID(rotation.NewToken.FamilyID),
		ExpiresAt: db.PGTimestamptz(rotation.NewToken.ExpiresAt),
	}); err != nil {
		return err
	}

	rowsAffected, err := qtx.RevokeRefreshTokenForRotation(ctx, dbgen.RevokeRefreshTokenForRotationParams{
		ID:                db.PGUUID(rotation.CurrentTokenID),
		RevokedAt:         db.PGTimestamptz(time.Now().UTC()),
		ReplacedByTokenID: db.PGUUID(rotation.ReplacedByID),
	})
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrInvalidRefresh
	}

	return tx.Commit(ctx)
}

func (r *postgresRepository) RevokeRefreshTokenFamily(ctx context.Context, familyID uuid.UUID) error {
	return r.queries.RevokeRefreshTokenFamily(ctx, dbgen.RevokeRefreshTokenFamilyParams{
		FamilyID:  db.PGUUID(familyID),
		RevokedAt: db.PGTimestamptz(time.Now().UTC()),
	})
}

func mapCreateUserRow(row dbgen.CreateUserRow) *User {
	return &User{
		ID:           db.UUID(row.ID),
		Email:        row.Email,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    db.Time(row.CreatedAt),
		UpdatedAt:    db.Time(row.UpdatedAt),
	}
}

func mapGetUserByEmailOrUsernameRow(row dbgen.GetUserByEmailOrUsernameRow) *User {
	return &User{
		ID:           db.UUID(row.ID),
		Email:        row.Email,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    db.Time(row.CreatedAt),
		UpdatedAt:    db.Time(row.UpdatedAt),
	}
}

func mapGetUserByIDRow(row dbgen.GetUserByIDRow) *User {
	return &User{
		ID:           db.UUID(row.ID),
		Email:        row.Email,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    db.Time(row.CreatedAt),
		UpdatedAt:    db.Time(row.UpdatedAt),
	}
}

func mapRefreshTokenRow(row dbgen.GetRefreshTokenRow) *RefreshToken {
	token := &RefreshToken{
		ID:        db.UUID(row.ID),
		UserID:    db.UUID(row.UserID),
		TokenHash: row.TokenHash,
		FamilyID:  db.UUID(row.FamilyID),
		ExpiresAt: db.Time(row.ExpiresAt),
		CreatedAt: db.Time(row.CreatedAt),
	}
	token.RevokedAt = db.TimePtr(row.RevokedAt)
	return token
}
