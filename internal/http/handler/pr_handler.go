package handler

import (
	"encoding/json"
	"net/http"
	"pr-reviewer-service/internal/models"
	"pr-reviewer-service/internal/service"
)

type PRHandler struct {
	service *service.PRService
}

func NewPRHandler(s *service.PRService) *PRHandler {
	return &PRHandler{service: s}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var pr models.PullRequest
	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreatePR(&pr); err != nil {
		if err.Error() == "PR_EXISTS" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"pr": pr})
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr, err := h.service.MergePR(req.PullRequestID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"pr": pr})
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr, newUserID, err := h.service.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		code := "NOT_FOUND"
		msg := "pr or user not found"

		if err.Error() == "PR_MERGED" {
			code = "PR_MERGED"
			msg = "cannot reassign on merged PR"
			w.WriteHeader(http.StatusConflict)
		} else if err.Error() == "NOT_ASSIGNED" {
			code = "NOT_ASSIGNED"
			msg = "reviewer is not assigned to this PR"
			w.WriteHeader(http.StatusConflict)
		} else if err.Error() == "NO_CANDIDATE" {
			code = "NO_CANDIDATE"
			msg = "no active replacement candidate in team"
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{"code": code, "message": msg},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr":          pr,
		"replaced_by": newUserID,
	})
}
