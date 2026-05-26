package reviews

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrReviewNotFound = errors.New("review not found")
	ErrInvalidReview  = errors.New("invalid review data")
	ErrReviewExists   = errors.New("review already exists")
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Create(ctx context.Context, params CreateReviewParams) (*Review, error) {
	if params.UserID == uuid.Nil || params.SourceID == uuid.Nil || !validRating(params.Rating) {
		return nil, ErrInvalidReview
	}

	review := &Review{
		UserID:   params.UserID,
		SourceID: params.SourceID,
		Rating:   params.Rating,
		Content:  params.Content,
		IsPublic: params.IsPublic,
	}
	created, err := s.repo.Create(ctx, review)
	if err != nil {
		s.logger.Error("failed to create review", "error", err)
		return nil, err
	}

	s.logger.Info("review created", "id", created.ID, "user_id", created.UserID, "source_id", created.SourceID)
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Review, error) {
	review, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get review", "error", err, "id", id)
		return nil, err
	}
	if review == nil {
		return nil, ErrReviewNotFound
	}
	return review, nil
}

func (s *Service) ListPublicByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Review, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListPublicByUser(ctx, userID, limit, offset)
}

func (s *Service) ListPublicBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Review, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListPublicBySource(ctx, sourceID, limit, offset)
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Review, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *Service) ListBySource(ctx context.Context, sourceID uuid.UUID, limit, offset int) ([]*Review, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListBySource(ctx, sourceID, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, params UpdateReviewParams) (*Review, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrReviewNotFound
	}

	if params.Rating != nil {
		if !validRating(*params.Rating) {
			return nil, ErrInvalidReview
		}
		existing.Rating = *params.Rating
	}
	if params.Content != nil {
		existing.Content = params.Content
	}
	if params.IsPublic != nil {
		existing.IsPublic = *params.IsPublic
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		s.logger.Error("failed to update review", "error", err, "id", id)
		return nil, err
	}

	s.logger.Info("review updated", "id", id)
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete review", "error", err, "id", id)
		return err
	}

	s.logger.Info("review deleted", "id", id)
	return nil
}

func validRating(rating int) bool {
	return rating >= 1 && rating <= 5
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
