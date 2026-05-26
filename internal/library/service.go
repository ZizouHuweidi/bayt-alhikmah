package library

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrItemNotFound = errors.New("library item not found")
	ErrInvalidItem  = errors.New("invalid library item")
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Create(ctx context.Context, params CreateItemParams) (*Item, error) {
	if params.UserID == uuid.Nil || params.SourceID == uuid.Nil || !validStatus(params.Status) {
		return nil, ErrInvalidItem
	}
	if params.Visibility == "" {
		params.Visibility = VisibilityPrivate
	}
	if !validVisibility(params.Visibility) || !validProgress(params.ProgressValue, params.ProgressUnit) {
		return nil, ErrInvalidItem
	}

	item := &Item{
		UserID:        params.UserID,
		SourceID:      params.SourceID,
		Status:        params.Status,
		ProgressValue: params.ProgressValue,
		ProgressUnit:  params.ProgressUnit,
		Visibility:    params.Visibility,
		StartedAt:     params.StartedAt,
		CompletedAt:   params.CompletedAt,
	}

	created, err := s.repo.Create(ctx, item)
	if err != nil {
		s.logger.Error("failed to create library item", "error", err)
		return nil, err
	}
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}
	return item, nil
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *Service) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Item, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListPublicByUser(ctx, userID, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, params UpdateItemParams) (*Item, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrItemNotFound
	}

	if params.Status != nil {
		if !validStatus(*params.Status) {
			return nil, ErrInvalidItem
		}
		existing.Status = *params.Status
	}
	if params.Visibility != nil {
		if !validVisibility(*params.Visibility) {
			return nil, ErrInvalidItem
		}
		existing.Visibility = *params.Visibility
	}
	if params.ProgressValue != nil {
		existing.ProgressValue = params.ProgressValue
	}
	if params.ProgressUnit != nil {
		existing.ProgressUnit = params.ProgressUnit
	}
	if !validProgress(existing.ProgressValue, existing.ProgressUnit) {
		return nil, ErrInvalidItem
	}
	if params.StartedAt != nil {
		existing.StartedAt = params.StartedAt
	}
	if params.CompletedAt != nil {
		existing.CompletedAt = params.CompletedAt
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		s.logger.Error("failed to update library item", "error", err, "id", id)
		return nil, err
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func validStatus(status Status) bool {
	switch status {
	case StatusToConsume, StatusInProgress, StatusCompleted, StatusPaused, StatusAbandoned:
		return true
	default:
		return false
	}
}

func validVisibility(visibility Visibility) bool {
	switch visibility {
	case VisibilityPrivate, VisibilityUnlisted, VisibilityPublic:
		return true
	default:
		return false
	}
}

func validProgress(value *int, unit *ProgressUnit) bool {
	if value != nil && *value < 0 {
		return false
	}
	if unit == nil {
		return true
	}
	switch *unit {
	case ProgressUnitPage, ProgressUnitPercent, ProgressUnitMinute, ProgressUnitSecond, ProgressUnitEpisode:
		return true
	default:
		return false
	}
}

func normalizePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
