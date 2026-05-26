package reviews

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Review struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	SourceID  uuid.UUID `json:"source_id"`
	Rating    int       `json:"rating"`
	Content   *string   `json:"content,omitempty"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, review *Review) (*Review, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Review, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Review, error)
	ListBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Review, error)
	ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Review, error)
	ListPublicBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Review, error)
	Update(ctx context.Context, review *Review) (*Review, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateReviewParams struct {
	UserID   uuid.UUID
	SourceID uuid.UUID
	Rating   int
	Content  *string
	IsPublic bool
}

type UpdateReviewParams struct {
	Rating   *int
	Content  *string
	IsPublic *bool
}
