package sources

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

var (
	ErrSourceNotFound = errors.New("source not found")
	ErrInvalidSource  = errors.New("invalid source data")
)

// Service provides business logic for sources
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new source service
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new source
func (s *Service) Create(ctx context.Context, params CreateSourceParams) (*Source, error) {
	if params.Title == "" {
		return nil, ErrInvalidSource
	}

	source := &Source{
		Title:       params.Title,
		Subtitle:    params.Subtitle,
		Type:        params.Type,
		Description: params.Description,
		AuthorID:    params.AuthorID,
		Publisher:   params.Publisher,
		ISBN:        params.ISBN,
		DOI:         params.DOI,
		URL:         params.URL,
		ExternalID:  params.ExternalID,
		Tags:        params.Tags,
		PublishedAt: params.PublishedAt,
	}

	created, err := s.repo.Create(ctx, source)
	if err != nil {
		s.logger.Error("failed to create source", "error", err)
		return nil, err
	}

	s.logger.Info("source created", "id", created.ID, "title", created.Title)
	return created, nil
}

// GetByID retrieves a source by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Source, error) {
	source, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get source", "error", err, "id", id)
		return nil, err
	}

	if source == nil {
		return nil, ErrSourceNotFound
	}

	return source, nil
}

// List retrieves a paginated list of sources
func (s *Service) List(ctx context.Context, limit, offset int) ([]*Source, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}

// ListByType retrieves sources filtered by type
func (s *Service) ListByType(ctx context.Context, sourceType SourceType, limit, offset int) ([]*Source, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListByType(ctx, sourceType, limit, offset)
}

// Update updates an existing source
func (s *Service) Update(ctx context.Context, id uuid.UUID, params UpdateSourceParams) (*Source, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrSourceNotFound
	}

	// Apply updates
	if params.Title != nil {
		existing.Title = *params.Title
	}
	if params.Subtitle != nil {
		existing.Subtitle = params.Subtitle
	}
	if params.Type != nil {
		existing.Type = *params.Type
	}
	if params.Description != nil {
		existing.Description = params.Description
	}
	if params.AuthorID != nil {
		existing.AuthorID = params.AuthorID
	}
	if params.Publisher != nil {
		existing.Publisher = params.Publisher
	}
	if params.ISBN != nil {
		existing.ISBN = params.ISBN
	}
	if params.DOI != nil {
		existing.DOI = params.DOI
	}
	if params.URL != nil {
		existing.URL = params.URL
	}
	if params.ExternalID != nil {
		existing.ExternalID = params.ExternalID
	}
	if params.Tags != nil {
		existing.Tags = params.Tags
	}
	if params.PublishedAt != nil {
		existing.PublishedAt = params.PublishedAt
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		s.logger.Error("failed to update source", "error", err, "id", id)
		return nil, err
	}

	s.logger.Info("source updated", "id", id)
	return updated, nil
}

// Delete removes a source
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete source", "error", err, "id", id)
		return err
	}

	s.logger.Info("source deleted", "id", id)
	return nil
}

// Search searches sources by title
func (s *Service) Search(ctx context.Context, query string, limit, offset int) ([]*Source, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.Search(ctx, query, limit, offset)
}
