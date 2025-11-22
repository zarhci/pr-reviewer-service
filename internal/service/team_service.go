package service

import (
	"errors"
	"pr-reviewer-service/internal/models"
	"pr-reviewer-service/internal/repository/memory"
)

var ErrTeamExists = errors.New("TEAM_EXISTS")
var ErrTeamNotFound = errors.New("NOT_FOUND")

type TeamService struct {
	repo     *memory.TeamRepository
	userRepo *memory.UserRepository
}

func NewTeamService(repo *memory.TeamRepository, userRepo *memory.UserRepository) *TeamService {
	return &TeamService{repo: repo, userRepo: userRepo}
}

func (s *TeamService) CreateOrUpdate(team *models.Team) error {
	_, exists := s.repo.Get(team.TeamName)
	if exists {

		s.repo.Update(team)
	} else {
		err := s.repo.Create(team)
		if err != nil {
			return ErrTeamExists
		}
	}

	for _, member := range team.Members {
		u := &models.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: team.TeamName,
			IsActive: member.IsActive,
		}
		s.userRepo.CreateOrUpdate(u)
	}

	return nil
}

func (s *TeamService) GetTeam(name string) (*models.Team, error) {
	team, exists := s.repo.Get(name)
	if !exists {
		return nil, ErrTeamNotFound
	}
	return team, nil
}
