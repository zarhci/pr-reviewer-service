package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"pr-reviewer-service/internal/models"
	"pr-reviewer-service/internal/repository"
)

var (
	ErrPRExists            = errors.New("PR_EXISTS")
	ErrPRNotFound          = errors.New("NOT_FOUND")
	ErrPRMerged            = errors.New("PR_MERGED")
	ErrReviewerNotAssigned = errors.New("NOT_ASSIGNED")
	ErrNoCandidate         = errors.New("NO_CANDIDATE")
)

type PRService struct {
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPRService(
	prRepo repository.PullRequestRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *PRService) CreatePR(
	ctx context.Context,
	pr *models.PullRequest,
) error {

	if _, err := s.prRepo.GetByID(ctx, pr.PullRequestID); err == nil {
		return ErrPRExists
	}

	author, err := s.userRepo.GetByID(ctx, pr.AuthorID)
	if err != nil {
		return ErrPRNotFound
	}

	team, err := s.teamRepo.GetByName(ctx, author.TeamName)
	if err != nil {
		return ErrPRNotFound
	}

	var candidates []string
	for _, m := range team.Members {
		if m.UserID != author.UserID && m.IsActive {
			candidates = append(candidates, m.UserID)
		}
	}

	if len(candidates) == 0 {
		return ErrNoCandidate
	}

	rand.Seed(time.Now().UnixNano())
	var assigned []string
	for i := 0; i < 2 && len(candidates) > 0; i++ {
		idx := rand.Intn(len(candidates))
		assigned = append(assigned, candidates[idx])
		candidates = append(candidates[:idx], candidates[idx+1:]...)
	}

	pr.AssignedReviewers = assigned
	now := time.Now()
	pr.CreatedAt = &now
	pr.Status = models.PRStatusOpen

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return ErrPRExists
	}

	return nil
}

func (s *PRService) MergePR(
	ctx context.Context,
	prID string,
) (*models.PullRequest, error) {

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == models.PRStatusMerged {
		return pr, nil
	}

	now := time.Now()
	pr.Status = models.PRStatusMerged
	pr.MergedAt = &now

	if err := s.prRepo.Update(ctx, pr); err != nil {
		return nil, ErrPRNotFound
	}

	return pr, nil
}

func (s *PRService) ReassignReviewer(
	ctx context.Context,
	prID string,
	oldUserID string,
) (*models.PullRequest, string, error) {

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
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

	oldUser, err := s.userRepo.GetByID(ctx, oldUserID)
	if err != nil {
		return nil, "", ErrPRNotFound
	}

	team, err := s.teamRepo.GetByName(ctx, oldUser.TeamName)
	if err != nil {
		return nil, "", ErrNoCandidate
	}

	var candidates []string
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

	if err := s.prRepo.Update(ctx, pr); err != nil {
		return nil, "", ErrPRNotFound
	}

	return pr, newUserID, nil
}
