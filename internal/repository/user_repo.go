package repository

import (
	"context"
	"database/sql"
	"errors"

	"pr-reviewer-service/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID string) (*models.User, error)
	SetActive(ctx context.Context, userID string, active bool) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	const q = `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		q,
		user.UserID,
		user.Username,
		user.TeamName,
		user.IsActive,
	)

	if err != nil {
		// unique violation → already exists
		return ErrExists
	}

	return nil
}

func (r *userRepository) GetByID(
	ctx context.Context,
	userID string,
) (*models.User, error) {
	const q = `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`

	var u models.User
	err := r.db.QueryRowContext(ctx, q, userID).
		Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) SetActive(
	ctx context.Context,
	userID string,
	active bool,
) (*models.User, error) {
	const q = `
		UPDATE users
		SET is_active = $2
		WHERE user_id = $1
		RETURNING user_id, username, team_name, is_active
	`

	var u models.User
	err := r.db.QueryRowContext(ctx, q, userID, active).
		Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}
