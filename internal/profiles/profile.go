package profiles

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Profile struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Username      string    `json:"username,omitempty"`
	DisplayName   *string   `json:"display_name,omitempty"`
	Bio           *string   `json:"bio,omitempty"`
	PublicProfile bool      `json:"public_profile"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Repository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Profile, error)
	GetPublicByUsername(ctx context.Context, username string) (*Profile, error)
	Upsert(ctx context.Context, profile *Profile) (*Profile, error)
}

type UpdateProfileParams struct {
	UserID        uuid.UUID
	DisplayName   *string
	Bio           *string
	PublicProfile *bool
}
