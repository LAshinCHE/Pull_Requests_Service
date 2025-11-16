package http

import (
	"encoding/json"
	"net/http"
	"pull_request_service/internal/domain"
)

type setUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *HTTPHandler) SetUserIsActive(w http.ResponseWriter, r *http.Request) {
	var req setUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "invalid JSON")
		return
	}

	user, err := h.User.SetUserActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		if err == domain.ErrNotFound {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"user": user})
}

func (h *HTTPHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "BAD_INPUT", "user_id required")
		return
	}

	prs, err := h.PullRequest.GetPRsByReviewer(r.Context(), userID)
	if err != nil {
		if err == domain.ErrNotFound {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
