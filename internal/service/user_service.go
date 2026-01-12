package service

import (
	"context"
	"errors"

	"pr-reviewer-service/internal/models"
	"pr-reviewer-service/internal/repository"
)

var ErrUserNotFound = errors.New("NOT_FOUND")

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SetActive(
	ctx context.Context,
	userID string,
	active bool,
) (*models.User, error) {
	u, err := s.repo.SetActive(ctx, userID, active)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func (s *UserService) GetByID(
	ctx context.Context,
	userID string,
) (*models.User, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}
