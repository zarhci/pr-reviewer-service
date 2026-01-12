package repository

import (
	"context"
	"database/sql"
	"errors"

	"pr-reviewer-service/internal/models"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr *models.PullRequest) error
	GetByID(ctx context.Context, prID string) (*models.PullRequest, error)
	Update(ctx context.Context, pr *models.PullRequest) error
	ListByTeam(ctx context.Context, teamName string) ([]models.PullRequestShort, error)
}

type pullRequestRepository struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) PullRequestRepository {
	return &pullRequestRepository{db: db}
}

func (r *pullRequestRepository) Create(
	ctx context.Context,
	pr *models.PullRequest,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const prQuery = `
		INSERT INTO pull_requests (
			pull_request_id,
			pull_request_name,
			author_id,
			status,
			created_at,
			merged_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = tx.ExecContext(
		ctx,
		prQuery,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
		pr.CreatedAt,
		pr.MergedAt,
	)
	if err != nil {
		return ErrPRExists
	}

	if len(pr.AssignedReviewers) > 0 {
		const reviewersQuery = `
			INSERT INTO pull_request_reviewers (pull_request_id, user_id)
			VALUES ($1, $2)
		`

		for _, userID := range pr.AssignedReviewers {
			if _, err := tx.ExecContext(
				ctx,
				reviewersQuery,
				pr.PullRequestID,
				userID,
			); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *pullRequestRepository) GetByID(
	ctx context.Context,
	prID string,
) (*models.PullRequest, error) {
	const prQuery = `
		SELECT pull_request_id,
		       pull_request_name,
		       author_id,
		       status,
		       created_at,
		       merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	var pr models.PullRequest
	err := r.db.QueryRowContext(ctx, prQuery, prID).
		Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
			&pr.CreatedAt,
			&pr.MergedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPRNotFound
		}
		return nil, err
	}

	const reviewersQuery = `
		SELECT user_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1
	`

	rows, err := r.db.QueryContext(ctx, reviewersQuery, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *pullRequestRepository) Update(
	ctx context.Context,
	pr *models.PullRequest,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const updatePR = `
		UPDATE pull_requests
		SET status = $2,
		    merged_at = $3
		WHERE pull_request_id = $1
	`

	res, err := tx.ExecContext(
		ctx,
		updatePR,
		pr.PullRequestID,
		pr.Status,
		pr.MergedAt,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrPRNotFound
	}

	// пересобираем reviewers
	const deleteReviewers = `
		DELETE FROM pull_request_reviewers
		WHERE pull_request_id = $1
	`
	if _, err := tx.ExecContext(ctx, deleteReviewers, pr.PullRequestID); err != nil {
		return err
	}

	const insertReviewer = `
		INSERT INTO pull_request_reviewers (pull_request_id, user_id)
		VALUES ($1, $2)
	`
	for _, userID := range pr.AssignedReviewers {
		if _, err := tx.ExecContext(
			ctx,
			insertReviewer,
			pr.PullRequestID,
			userID,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *pullRequestRepository) ListByTeam(
	ctx context.Context,
	teamName string,
) ([]models.PullRequestShort, error) {
	const q = `
		SELECT pr.pull_request_id,
		       pr.pull_request_name,
		       pr.author_id,
		       pr.status
		FROM pull_requests pr
		JOIN users u ON u.user_id = pr.author_id
		WHERE u.team_name = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, q, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
		); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}
