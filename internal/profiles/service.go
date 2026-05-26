package profiles

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
	ErrInvalidProfile  = errors.New("invalid profile data")
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) GetOwn(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	profile, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get profile", "error", err, "user_id", userID)
		return nil, err
	}
	if profile == nil {
		return s.repo.Upsert(ctx, &Profile{UserID: userID})
	}
	return profile, nil
}

func (s *Service) GetPublicByUsername(ctx context.Context, username string) (*Profile, error) {
	if username == "" {
		return nil, ErrInvalidProfile
	}
	profile, err := s.repo.GetPublicByUsername(ctx, username)
	if err != nil {
		s.logger.Error("failed to get public profile", "error", err, "username", username)
		return nil, err
	}
	if profile == nil {
		return nil, ErrProfileNotFound
	}
	return profile, nil
}

func (s *Service) Update(ctx context.Context, params UpdateProfileParams) (*Profile, error) {
	if params.UserID == uuid.Nil {
		return nil, ErrInvalidProfile
	}

	existing, err := s.repo.GetByUserID(ctx, params.UserID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		existing = &Profile{UserID: params.UserID}
	}

	if params.DisplayName != nil {
		existing.DisplayName = params.DisplayName
	}
	if params.Bio != nil {
		existing.Bio = params.Bio
	}
	if params.PublicProfile != nil {
		existing.PublicProfile = *params.PublicProfile
	}

	updated, err := s.repo.Upsert(ctx, existing)
	if err != nil {
		s.logger.Error("failed to update profile", "error", err, "user_id", params.UserID)
		return nil, err
	}

	s.logger.Info("profile updated", "user_id", params.UserID)
	return updated, nil
}
