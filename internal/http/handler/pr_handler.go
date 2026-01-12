package handler

import (
	"net/http"

	"github.com/labstack/echo"

	"pr-reviewer-service/internal/models"
	"pr-reviewer-service/internal/service"
)

type PRHandler struct {
	prService *service.PRService
}

func NewPRHandler(prService *service.PRService) *PRHandler {
	return &PRHandler{
		prService: prService,
	}
}

/*
POST /pullRequest/create

	{
	  "pull_request_id": "pr1",
	  "pull_request_name": "Add feature",
	  "author_id": "u1"
	}
*/
type createPRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

func (h *PRHandler) CreatePR(c echo.Context) error {
	ctx := c.Request().Context()

	var req createPRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	if req.PullRequestID == "" || req.AuthorID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "VALIDATION_ERROR",
				"message": "pull_request_id and author_id are required",
			},
		})
	}

	pr := &models.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}

	if err := h.prService.CreatePR(ctx, pr); err != nil {
		switch err {
		case service.ErrPRExists:
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"error": map[string]string{
					"code":    "PR_EXISTS",
					"message": "pull request already exists",
				},
			})
		case service.ErrNoCandidate:
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"error": map[string]string{
					"code":    "NO_CANDIDATE",
					"message": "no active reviewers available",
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"code":    "INTERNAL_ERROR",
					"message": "could not create pull request",
				},
			})
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"status": "ok",
	})
}

/*
POST /pullRequest/merge

	{
	  "pull_request_id": "pr1"
	}
*/
type mergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

func (h *PRHandler) MergePR(c echo.Context) error {
	ctx := c.Request().Context()

	var req mergePRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	pr, err := h.prService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		switch err {
		case service.ErrPRNotFound:
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{
					"code":    "PR_NOT_FOUND",
					"message": "pull request not found",
				},
			})
		case service.ErrPRMerged:
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"error": map[string]string{
					"code":    "PR_ALREADY_MERGED",
					"message": "pull request already merged",
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"code":    "INTERNAL_ERROR",
					"message": "could not merge pull request",
				},
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"pull_request": pr,
	})
}

/*
POST /pullRequest/reassign

	{
	  "pull_request_id": "pr1",
	  "old_user_id": "u2"
	}
*/
type reassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

func (h *PRHandler) ReassignReviewer(c echo.Context) error {
	ctx := c.Request().Context()

	var req reassignReviewerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	pr, newUserID, err := h.prService.ReassignReviewer(
		ctx,
		req.PullRequestID,
		req.OldUserID,
	)
	if err != nil {
		switch err {
		case service.ErrReviewerNotAssigned:
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": map[string]string{
					"code":    "REVIEWER_NOT_ASSIGNED",
					"message": "reviewer not assigned to PR",
				},
			})
		case service.ErrNoCandidate:
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"error": map[string]string{
					"code":    "NO_CANDIDATE",
					"message": "no replacement reviewer available",
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"code":    "INTERNAL_ERROR",
					"message": "could not reassign reviewer",
				},
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"pull_request": pr,
		"new_user_id":  newUserID,
	})
}
