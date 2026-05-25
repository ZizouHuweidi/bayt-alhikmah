package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSignup      = errors.New("invalid signup data")
	ErrInvalidRefresh     = errors.New("invalid refresh token")
	ErrUserNotFound       = errors.New("user not found")
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9_][a-z0-9_-]{2,31}$`)

type Service struct {
	repo       Repository
	tokens     *TokenManager
	refreshTTL time.Duration
	logger     *slog.Logger
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

func NewService(repo Repository, tokens *TokenManager, refreshTTL time.Duration, logger *slog.Logger) *Service {
	return &Service{repo: repo, tokens: tokens, refreshTTL: refreshTTL, logger: logger}
}

func (s *Service) Register(ctx context.Context, email, username, password string) (*User, AuthTokens, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.ToLower(strings.TrimSpace(username))
	if !validEmail(email) || !usernamePattern.MatchString(username) || len(password) < 12 {
		return nil, AuthTokens{}, ErrInvalidSignup
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, AuthTokens{}, err
	}
	userID, err := uuid.NewV7()
	if err != nil {
		return nil, AuthTokens{}, err
	}

	user, err := s.repo.CreateUser(ctx, User{ID: userID, Email: email, Username: username, PasswordHash: passwordHash})
	if err != nil {
		return nil, AuthTokens{}, err
	}

	tokens, err := s.issueTokens(ctx, *user)
	if err != nil {
		return nil, AuthTokens{}, err
	}
	return user, tokens, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (*User, AuthTokens, error) {
	user, err := s.repo.GetUserByEmailOrUsername(ctx, strings.TrimSpace(login))
	if err != nil {
		return nil, AuthTokens{}, err
	}
	if user == nil {
		return nil, AuthTokens{}, ErrInvalidCredentials
	}

	valid, err := VerifyPassword(password, user.PasswordHash)
	if err != nil || !valid {
		return nil, AuthTokens{}, ErrInvalidCredentials
	}

	tokens, err := s.issueTokens(ctx, *user)
	if err != nil {
		return nil, AuthTokens{}, err
	}
	return user, tokens, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (AuthTokens, error) {
	existing, err := s.repo.GetRefreshToken(ctx, HashRefreshToken(refreshToken))
	if err != nil {
		return AuthTokens{}, err
	}
	if existing == nil || time.Now().UTC().After(existing.ExpiresAt) {
		return AuthTokens{}, ErrInvalidRefresh
	}
	if existing.RevokedAt != nil {
		_ = s.repo.RevokeRefreshTokenFamily(ctx, existing.FamilyID)
		return AuthTokens{}, ErrInvalidRefresh
	}

	user, err := s.repo.GetUserByID(ctx, existing.UserID)
	if err != nil {
		return AuthTokens{}, err
	}
	if user == nil {
		return AuthTokens{}, ErrUserNotFound
	}

	accessToken, err := s.tokens.CreateAccessToken(*user)
	if err != nil {
		return AuthTokens{}, err
	}
	newRefresh, newHash, err := NewRefreshToken()
	if err != nil {
		return AuthTokens{}, err
	}
	newID, err := uuid.NewV7()
	if err != nil {
		return AuthTokens{}, err
	}
	newToken := RefreshToken{
		ID:        newID,
		UserID:    user.ID,
		TokenHash: newHash,
		FamilyID:  existing.FamilyID,
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL),
	}
	if err := s.repo.RotateRefreshToken(ctx, RefreshTokenRotation{
		CurrentTokenID: existing.ID,
		NewToken:       newToken,
		ReplacedByID:   newID,
	}); err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{AccessToken: accessToken, RefreshToken: newRefresh, TokenType: "Bearer", ExpiresIn: int64(s.tokens.accessTTL.Seconds())}, nil
}

func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *Service) VerifyAccessToken(rawToken string) (*AccessClaims, error) {
	return s.tokens.VerifyAccessToken(rawToken)
}

func (s *Service) issueTokens(ctx context.Context, user User) (AuthTokens, error) {
	accessToken, err := s.tokens.CreateAccessToken(user)
	if err != nil {
		return AuthTokens{}, err
	}
	refreshToken, refreshHash, err := NewRefreshToken()
	if err != nil {
		return AuthTokens{}, err
	}
	refreshID, err := uuid.NewV7()
	if err != nil {
		return AuthTokens{}, err
	}
	familyID, err := uuid.NewV7()
	if err != nil {
		return AuthTokens{}, err
	}

	if err := s.repo.CreateRefreshToken(ctx, RefreshToken{
		ID:        refreshID,
		UserID:    user.ID,
		TokenHash: refreshHash,
		FamilyID:  familyID,
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL),
	}); err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "Bearer", ExpiresIn: int64(s.tokens.accessTTL.Seconds())}, nil
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
