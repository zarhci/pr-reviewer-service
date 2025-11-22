package service

import (
	"errors"
	"math/rand"
	"pr-reviewer-service/internal/models"
	"pr-reviewer-service/internal/repository/memory"
	"time"
)

var (
	ErrPRExists            = errors.New("PR_EXISTS")
	ErrPRNotFound          = errors.New("NOT_FOUND")
	ErrPRMerged            = errors.New("PR_MERGED")
	ErrReviewerNotAssigned = errors.New("NOT_ASSIGNED")
	ErrNoCandidate         = errors.New("NO_CANDIDATE")
)

type PRService struct {
	prRepo   *memory.PRRepository
	userRepo *memory.UserRepository
	teamRepo *memory.TeamRepository
}

func NewPRService(prRepo *memory.PRRepository, userRepo *memory.UserRepository, teamRepo *memory.TeamRepository) *PRService {
	return &PRService{prRepo: prRepo, userRepo: userRepo, teamRepo: teamRepo}
}

func (s *PRService) CreatePR(pr *models.PullRequest) error {

	if _, exists := s.prRepo.Get(pr.PullRequestID); exists {
		return ErrPRExists
	}

	author, ok := s.userRepo.GetByID(pr.AuthorID)
	if !ok {
		return ErrPRNotFound
	}

	team, ok := s.teamRepo.Get(author.TeamName)
	if !ok {
		return ErrPRNotFound
	}

	candidates := []string{}
	for _, m := range team.Members {
		if m.UserID != author.UserID && m.IsActive {
			candidates = append(candidates, m.UserID)
		}
	}

	assigned := []string{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 2 && len(candidates) > 0; i++ {
		idx := rand.Intn(len(candidates))
		assigned = append(assigned, candidates[idx])

		candidates = append(candidates[:idx], candidates[idx+1:]...)
	}

	pr.AssignedReviewers = assigned
	now := time.Now()
	pr.CreatedAt = &now
	pr.Status = models.PRStatusOpen

	return s.prRepo.Create(pr)
}

func (s *PRService) MergePR(prID string) (*models.PullRequest, error) {
	pr, ok := s.prRepo.Get(prID)
	if !ok {
		return nil, ErrPRNotFound
	}

	if pr.Status == models.PRStatusMerged {
		return pr, nil
	}

	now := time.Now()
	pr.Status = models.PRStatusMerged
	pr.MergedAt = &now

	s.prRepo.Update(pr)
	return pr, nil
}

func (s *PRService) ReassignReviewer(prID string, oldUserID string) (*models.PullRequest, string, error) {
	pr, ok := s.prRepo.Get(prID)
	if !ok {
		return nil, "", ErrPRNotFound
	}

	if pr.Status == models.PRStatusMerged {
		return nil, "", ErrPRMerged
	}

	found := false
	for _, rev := range pr.AssignedReviewers {
		if rev == oldUserID {
			found = true
			break
		}
	}
	if !found {
		return nil, "", ErrReviewerNotAssigned
	}

	oldUser, ok := s.userRepo.GetByID(oldUserID)
	if !ok {
		return nil, "", ErrPRNotFound
	}

	team, ok := s.teamRepo.Get(oldUser.TeamName)
	if !ok {
		return nil, "", ErrNoCandidate
	}

	candidates := []string{}
	for _, m := range team.Members {
		if m.UserID != oldUserID && m.IsActive {
			candidates = append(candidates, m.UserID)
		}
	}

	if len(candidates) == 0 {
		return nil, "", ErrNoCandidate
	}

	rand.Seed(time.Now().UnixNano())
	newUserID := candidates[rand.Intn(len(candidates))]

	for i, rev := range pr.AssignedReviewers {
		if rev == oldUserID {
			pr.AssignedReviewers[i] = newUserID
			break
		}
	}

	s.prRepo.Update(pr)
	return pr, newUserID, nil
}

func (s *PRService) GetPRsForReviewer(userID string) []models.PullRequestShort {
	return s.prRepo.GetWhereReviewerIs(userID)
}
