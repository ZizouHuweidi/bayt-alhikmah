// Package notes implements the Notes bounded context for managing
// user notes, annotations, and highlights on knowledge sources.
package notes

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ContentType represents the type of note content
type ContentType string

const (
	ContentTypeNote       ContentType = "note"
	ContentTypeAnnotation ContentType = "annotation"
	ContentTypeQuote      ContentType = "quote"
	ContentTypeSummary    ContentType = "summary"
	ContentTypeReflection ContentType = "reflection"
)

// Note represents a user's thought or annotation on a source
type Note struct {
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"user_id"`
	SourceID    *uuid.UUID  `json:"source_id,omitempty"`
	Content     string      `json:"content"`
	ContentType ContentType `json:"content_type"`
	IsPublic    bool        `json:"is_public"`
	Annotations []string    `json:"annotations,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Repository defines the interface for note data access
type Repository interface {
	Create(ctx context.Context, note *Note) (*Note, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Note, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Note, error)
	ListBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Note, error)
	ListPublic(ctx context.Context, limit, offset int) ([]*Note, error)
	Update(ctx context.Context, note *Note) (*Note, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}

// CreateNoteParams contains parameters for creating a note
type CreateNoteParams struct {
	UserID      uuid.UUID
	SourceID    *uuid.UUID
	Content     string
	ContentType ContentType
	IsPublic    bool
	Annotations []string
	Tags        []string
}

// UpdateNoteParams contains parameters for updating a note
type UpdateNoteParams struct {
	Content     *string
	ContentType *ContentType
	IsPublic    *bool
	Annotations []string
	Tags        []string
}
