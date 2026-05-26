package library

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Status string

const (
	StatusToConsume  Status = "to_consume"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusPaused     Status = "paused"
	StatusAbandoned  Status = "abandoned"
)

type ProgressUnit string

const (
	ProgressUnitPage    ProgressUnit = "page"
	ProgressUnitPercent ProgressUnit = "percent"
	ProgressUnitMinute  ProgressUnit = "minute"
	ProgressUnitSecond  ProgressUnit = "second"
	ProgressUnitEpisode ProgressUnit = "episode"
)

type Visibility string

const (
	VisibilityPrivate  Visibility = "private"
	VisibilityUnlisted Visibility = "unlisted"
	VisibilityPublic   Visibility = "public"
)

type Item struct {
	ID            uuid.UUID     `json:"id"`
	UserID        uuid.UUID     `json:"user_id"`
	SourceID      uuid.UUID     `json:"source_id"`
	Status        Status        `json:"status"`
	ProgressValue *int          `json:"progress_value,omitempty"`
	ProgressUnit  *ProgressUnit `json:"progress_unit,omitempty"`
	Visibility    Visibility    `json:"visibility"`
	StartedAt     *time.Time    `json:"started_at,omitempty"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type SourceSummary struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Subtitle  *string   `json:"subtitle,omitempty"`
	Type      string    `json:"type"`
	Publisher *string   `json:"publisher,omitempty"`
	ISBN      *string   `json:"isbn,omitempty"`
}

type ItemWithSource struct {
	*Item
	Source *SourceSummary `json:"source"`
}

type Repository interface {
	Create(ctx context.Context, item *Item) (*Item, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Item, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error)
	ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error)
	ListPublicByUsername(ctx context.Context, username string, limit, offset int) ([]*Item, error)
	ListByUserWithSources(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ItemWithSource, error)
	ListPublicByUsernameWithSources(ctx context.Context, username string, limit, offset int) ([]*ItemWithSource, error)
	Update(ctx context.Context, item *Item) (*Item, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateItemParams struct {
	UserID        uuid.UUID
	SourceID      uuid.UUID
	Status        Status
	ProgressValue *int
	ProgressUnit  *ProgressUnit
	Visibility    Visibility
	StartedAt     *time.Time
	CompletedAt   *time.Time
}

type UpdateItemParams struct {
	Status        *Status
	ProgressValue *int
	ProgressUnit  *ProgressUnit
	Visibility    *Visibility
	StartedAt     *time.Time
	CompletedAt   *time.Time
}
