// Package sources implements the Sources bounded context for managing
// knowledge sources like books, papers, podcasts, videos, and articles.
package sources

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

// SourceType represents the type of knowledge source
type SourceType string

const (
	SourceTypeBook    SourceType = "book"
	SourceTypePaper   SourceType = "paper"
	SourceTypePodcast SourceType = "podcast"
	SourceTypeVideo   SourceType = "video"
	SourceTypeArticle SourceType = "article"
	SourceTypeEssay   SourceType = "essay"
)

// Source represents a knowledge source entity
type Source struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title" db:"title"`
	Subtitle    *string    `json:"subtitle,omitempty" db:"subtitle"`
	Type        SourceType `json:"type" db:"type"`
	Description *string    `json:"description,omitempty" db:"description"`
	AuthorID    *uuid.UUID `json:"author_id,omitempty" db:"author_id"`
	Publisher   *string    `json:"publisher,omitempty" db:"publisher"`
	ISBN        *string    `json:"isbn,omitempty" db:"isbn"`
	DOI         *string    `json:"doi,omitempty" db:"doi"`
	URL         *string    `json:"url,omitempty" db:"url"`
	ExternalID  *string    `json:"external_id,omitempty" db:"external_id"`
	Tags        []string   `json:"tags,omitempty" db:"-"`
	PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Author represents a source author
type Author struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Bio       *string   `json:"bio,omitempty"`
	BirthDate *string   `json:"birth_date,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Tag represents a taxonomy tag
type Tag struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Repository defines the interface for source data access
type Repository interface {
	Create(ctx context.Context, source *Source) (*Source, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Source, error)
	List(ctx context.Context, limit, offset int) ([]*Source, error)
	ListByType(ctx context.Context, sourceType SourceType, limit, offset int) ([]*Source, error)
	Update(ctx context.Context, source *Source) (*Source, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, limit, offset int) ([]*Source, error)
	Count(ctx context.Context) (int64, error)
}

// CreateSourceParams contains parameters for creating a source
type CreateSourceParams struct {
	Title       string
	Subtitle    *string
	Type        SourceType
	Description *string
	AuthorID    *uuid.UUID
	Publisher   *string
	ISBN        *string
	DOI         *string
	URL         *string
	ExternalID  *string
	Tags        []string
	PublishedAt *time.Time
}

// UpdateSourceParams contains parameters for updating a source
type UpdateSourceParams struct {
	Title       *string
	Subtitle    *string
	Type        *SourceType
	Description *string
	AuthorID    *uuid.UUID
	Publisher   *string
	ISBN        *string
	DOI         *string
	URL         *string
	ExternalID  *string
	Tags        []string
	PublishedAt *time.Time
}
