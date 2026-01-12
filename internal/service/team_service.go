package service

import (
	"context"
	"errors"

	"pr-reviewer-service/internal/models"
	"pr-reviewer-service/internal/repository"
)

var ErrTeamExists = errors.New("TEAM_EXISTS")
var ErrTeamNotFound = errors.New("NOT_FOUND")

type TeamService struct {
	repo     repository.TeamRepository
	userRepo repository.UserRepository
}

func NewTeamService(
	repo repository.TeamRepository,
	userRepo repository.UserRepository,
) *TeamService {
	return &TeamService{repo: repo, userRepo: userRepo}
}

func (s *TeamService) CreateTeam(
	ctx context.Context,
	team *models.Team,
) error {
	err := s.repo.Create(ctx, team)
	if err != nil {
		return ErrTeamExists
	}

	for _, member := range team.Members {
		u := &models.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: team.TeamName,
			IsActive: member.IsActive,
		}
		_ = s.userRepo.Create(ctx, u)
	}

	return nil
}

func (s *TeamService) GetTeam(
	ctx context.Context,
	name string,
) (*models.Team, error) {
	team, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}
