package memory

import "pr-reviewer-service/internal/models"

type UserRepository struct {
	store *Storage
}

func NewUserRepository(s *Storage) *UserRepository {
	return &UserRepository{store: s}
}

func (r *UserRepository) CreateOrUpdate(u *models.User) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	r.store.Users[u.UserID] = u
}

func (r *UserRepository) GetByID(id string) (*models.User, bool) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	u, ok := r.store.Users[id]
	return u, ok
}

func (r *UserRepository) SetActive(id string, active bool) (*models.User, bool) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	u, ok := r.store.Users[id]
	if !ok {
		return nil, false
	}
	u.IsActive = active
	return u, true
}
