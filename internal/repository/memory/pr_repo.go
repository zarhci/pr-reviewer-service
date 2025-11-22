package memory

import "pr-reviewer-service/internal/models"

type PRRepository struct {
	store *Storage
}

func NewPRRepository(s *Storage) *PRRepository {
	return &PRRepository{store: s}
}

func (r *PRRepository) Create(pr *models.PullRequest) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	if _, exists := r.store.PRs[pr.PullRequestID]; exists {
		return ErrPRExists
	}

	r.store.PRs[pr.PullRequestID] = pr
	return nil
}

func (r *PRRepository) Update(pr *models.PullRequest) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	r.store.PRs[pr.PullRequestID] = pr
}

func (r *PRRepository) Get(id string) (*models.PullRequest, bool) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	pr, ok := r.store.PRs[id]
	return pr, ok
}

func (r *PRRepository) GetWhereReviewerIs(userID string) []models.PullRequestShort {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	result := []models.PullRequestShort{}

	for _, pr := range r.store.PRs {
		for _, rev := range pr.AssignedReviewers {
			if rev == userID {
				result = append(result, models.PullRequestShort{
					PullRequestID:   pr.PullRequestID,
					PullRequestName: pr.PullRequestName,
					AuthorID:        pr.AuthorID,
					Status:          pr.Status,
				})
				break
			}
		}
	}

	return result
}
