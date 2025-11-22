package memory

import (
	"sync"

	"pr-reviewer-service/internal/models"
)

type Storage struct {
	mu sync.RWMutex

	Users map[string]*models.User
	Teams map[string]*models.Team
	PRs   map[string]*models.PullRequest
}

func NewStorage() *Storage {
	return &Storage{
		Users: make(map[string]*models.User),
		Teams: make(map[string]*models.Team),
		PRs:   make(map[string]*models.PullRequest),
	}
}
