// Package sources implements the Sources bounded context for managing
// knowledge sources like books, papers, podcasts, videos, and articles.
package sources

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	Title       string     `json:"title"`
	Subtitle    *string    `json:"subtitle,omitempty"`
	Type        SourceType `json:"type"`
	Description *string    `json:"description,omitempty"`
	AuthorID    *uuid.UUID `json:"author_id,omitempty"`
	Publisher   *string    `json:"publisher,omitempty"`
	ISBN        *string    `json:"isbn,omitempty"`
	DOI         *string    `json:"doi,omitempty"`
	URL         *string    `json:"url,omitempty"`
	ExternalID  *string    `json:"external_id,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
