package service

import (
	"errors"

	"pr-reviewer-service/internal/models"
	"pr-reviewer-service/internal/repository/memory"
)

var ErrUserNotFound = errors.New("NOT_FOUND")

type UserService struct {
	repo *memory.UserRepository
}

func NewUserService(repo *memory.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SetActive(userID string, active bool) (*models.User, error) {
	u, ok := s.repo.SetActive(userID, active)
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func (s *UserService) GetByID(userID string) (*models.User, error) {
	u, ok := s.repo.GetByID(userID)
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}
