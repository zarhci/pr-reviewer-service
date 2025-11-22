package memory

import "pr-reviewer-service/internal/models"

type TeamRepository struct {
	store *Storage
}

func NewTeamRepository(s *Storage) *TeamRepository {
	return &TeamRepository{store: s}
}

func (r *TeamRepository) Create(team *models.Team) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	if _, exists := r.store.Teams[team.TeamName]; exists {
		return ErrTeamExists
	}

	r.store.Teams[team.TeamName] = team
	return nil
}

func (r *TeamRepository) Update(team *models.Team) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	r.store.Teams[team.TeamName] = team
}

func (r *TeamRepository) Get(name string) (*models.Team, bool) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	t, ok := r.store.Teams[name]
	return t, ok
}
