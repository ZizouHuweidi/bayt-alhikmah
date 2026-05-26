package collections

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Collection struct {
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"user_id"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	IsPublic    bool        `json:"is_public"`
	SourceIDs   []uuid.UUID `json:"source_ids,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, collection *Collection) (*Collection, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Collection, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Collection, error)
	ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Collection, error)
	Update(ctx context.Context, collection *Collection) (*Collection, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateCollectionParams struct {
	UserID      uuid.UUID
	Name        string
	Description *string
	IsPublic    bool
	SourceIDs   []uuid.UUID
}

type UpdateCollectionParams struct {
	Name        *string
	Description *string
	IsPublic    *bool
	SourceIDs   []uuid.UUID
}
