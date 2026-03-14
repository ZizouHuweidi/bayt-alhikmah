package notes

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrInvalidNote  = errors.New("invalid note data")
)

// Service provides business logic for notes
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new note service
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new note
func (s *Service) Create(ctx context.Context, params CreateNoteParams) (*Note, error) {
	if params.Content == "" {
		return nil, ErrInvalidNote
	}

	note := &Note{
		UserID:      params.UserID,
		SourceID:    params.SourceID,
		Content:     params.Content,
		ContentType: params.ContentType,
		IsPublic:    params.IsPublic,
		Annotations: params.Annotations,
		Tags:        params.Tags,
	}

	created, err := s.repo.Create(ctx, note)
	if err != nil {
		s.logger.Error("failed to create note", "error", err)
		return nil, err
	}

	s.logger.Info("note created", "id", created.ID, "user_id", created.UserID)
	return created, nil
}

// GetByID retrieves a note by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Note, error) {
	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get note", "error", err, "id", id)
		return nil, err
	}

	if note == nil {
		return nil, ErrNoteNotFound
	}

	return note, nil
}

// ListByUser retrieves notes for a specific user
func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Note, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListByUser(ctx, userID, limit, offset)
}

// ListBySource retrieves notes for a specific source
func (s *Service) ListBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Note, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListBySource(ctx, sourceID, limit, offset)
}

// ListPublic retrieves public notes
func (s *Service) ListPublic(ctx context.Context, limit, offset int) ([]*Note, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListPublic(ctx, limit, offset)
}

// Update updates an existing note
func (s *Service) Update(ctx context.Context, id uuid.UUID, params UpdateNoteParams) (*Note, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNoteNotFound
	}

	// Apply updates
	if params.Content != nil {
		existing.Content = *params.Content
	}
	if params.ContentType != nil {
		existing.ContentType = *params.ContentType
	}
	if params.IsPublic != nil {
		existing.IsPublic = *params.IsPublic
	}
	if params.Annotations != nil {
		existing.Annotations = params.Annotations
	}
	if params.Tags != nil {
		existing.Tags = params.Tags
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		s.logger.Error("failed to update note", "error", err, "id", id)
		return nil, err
	}

	s.logger.Info("note updated", "id", id)
	return updated, nil
}

// Delete removes a note
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete note", "error", err, "id", id)
		return err
	}

	s.logger.Info("note deleted", "id", id)
	return nil
}
