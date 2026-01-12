package repository

import (
	"context"
	"database/sql"
	"errors"
	"pr-reviewer-service/internal/models"
)

type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	GetByName(ctx context.Context, teamName string) (*models.Team, error)
	AddUser(ctx context.Context, teamName, userID string) error
}

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team *models.Team) error {
	const q = `
		INSERT INTO teams (team_name)
		VALUES ($1)
	`

	_, err := r.db.ExecContext(ctx, q, team.TeamName)
	if err != nil {
		// уникальность PK → команда уже существует
		return ErrTeamExists
	}

	return nil
}

func (r *teamRepository) GetByName(
	ctx context.Context,
	teamName string,
) (*models.Team, error) {
	// 1. проверяем, что команда существует
	const teamQuery = `
		SELECT team_name
		FROM teams
		WHERE team_name = $1
	`

	var name string
	err := r.db.QueryRowContext(ctx, teamQuery, teamName).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}

	// 2. получаем участников команды
	const membersQuery = `
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
	`

	rows, err := r.db.QueryContext(ctx, membersQuery, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		var m models.TeamMember
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}

func (r *teamRepository) AddUser(
	ctx context.Context,
	teamName string,
	userID string,
) error {
	const q = `
		UPDATE users
		SET team_name = $1
		WHERE user_id = $2
	`

	res, err := r.db.ExecContext(ctx, q, teamName, userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		// либо пользователя нет, либо команды нет
		// но с точки зрения domain — это NOT_FOUND
		return ErrNotFound
	}

	return nil
}
