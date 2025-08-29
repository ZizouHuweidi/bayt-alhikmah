package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("not found")

type Service interface {
	Register(ctx context.Context, username, email, password string) (User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, username, email, password string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, username, email, string(hashedPassword))
	if err != nil {
		// TODO: Here you can check for specific DB errors, like a duplicate email,
		return User{}, err
	}

	return user, nil
}
