package handler

import (
	"encoding/json"
	"net/http"
	"pr-reviewer-service/internal/service"
)

type UserHandler struct {
	service   *service.UserService
	prService *service.PRService
}

func NewUserHandler(s *service.UserService, prService *service.PRService) *UserHandler {
	return &UserHandler{service: s, prService: prService}
}

func (h *UserHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.SetActive(req.UserID, req.IsActive)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    "NOT_FOUND",
				"message": "user not found",
			},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"user": user})
}

func (h *UserHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	prs := h.prService.GetPRsForReviewer(userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
