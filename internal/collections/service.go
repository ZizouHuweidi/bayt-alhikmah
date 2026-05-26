package collections

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrCollectionNotFound = errors.New("collection not found")
	ErrInvalidCollection  = errors.New("invalid collection data")
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Create(ctx context.Context, params CreateCollectionParams) (*Collection, error) {
	if params.UserID == uuid.Nil || params.Name == "" {
		return nil, ErrInvalidCollection
	}

	collection := &Collection{
		UserID:      params.UserID,
		Name:        params.Name,
		Description: params.Description,
		IsPublic:    params.IsPublic,
		SourceIDs:   params.SourceIDs,
	}
	created, err := s.repo.Create(ctx, collection)
	if err != nil {
		s.logger.Error("failed to create collection", "error", err)
		return nil, err
	}

	s.logger.Info("collection created", "id", created.ID, "user_id", created.UserID)
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Collection, error) {
	collection, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get collection", "error", err, "id", id)
		return nil, err
	}
	if collection == nil {
		return nil, ErrCollectionNotFound
	}
	return collection, nil
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Collection, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *Service) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Collection, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListPublicByUser(ctx, userID, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, params UpdateCollectionParams) (*Collection, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrCollectionNotFound
	}

	if params.Name != nil {
		if *params.Name == "" {
			return nil, ErrInvalidCollection
		}
		existing.Name = *params.Name
	}
	if params.Description != nil {
		existing.Description = params.Description
	}
	if params.IsPublic != nil {
		existing.IsPublic = *params.IsPublic
	}
	if params.SourceIDs != nil {
		existing.SourceIDs = params.SourceIDs
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		s.logger.Error("failed to update collection", "error", err, "id", id)
		return nil, err
	}

	s.logger.Info("collection updated", "id", id)
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete collection", "error", err, "id", id)
		return err
	}

	s.logger.Info("collection deleted", "id", id)
	return nil
}

func normalizePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
