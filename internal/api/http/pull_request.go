package http

import (
	"encoding/json"
	"net/http"
	"pull_request_service/internal/domain"
)

type createPRRequest struct {
	ID     string `json:"pull_request_id"`
	Name   string `json:"pull_request_name"`
	Author string `json:"author_id"`
}

type mergePRRequest struct {
	ID string `json:"pull_request_id"`
}

type reassignRequest struct {
	ID      string `json:"pull_request_id"`
	OldUser string `json:"old_user_id"`
}

func (h *HTTPHandler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req createPRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "invalid JSON")
		return
	}

	pr, err := h.PullRequest.Create(r.Context(), req.ID, req.Name, req.Author)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "author or team not found")
		case domain.ErrPRExists:
			writeError(w, http.StatusConflict, "PR_EXISTS", "PR id already exists")
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"pr": pr})
}

func (h *HTTPHandler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var req mergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "invalid JSON")
		return
	}

	pr, err := h.PullRequest.Merge(r.Context(), req.ID)
	if err != nil {
		if err == domain.ErrNotFound {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "PR not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"pr": pr})
}

func (h *HTTPHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req reassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "invalid JSON")
		return
	}

	pr, replacedBy, err := h.PullRequest.Reassign(r.Context(), req.ID, req.OldUser)
	if err != nil {
		switch err {
		case domain.ErrPRMerged:
			writeError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
		case domain.ErrNotAssigned:
			writeError(w, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
		case domain.ErrNoCandidate:
			writeError(w, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
		case domain.ErrNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "PR or user not found")
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pr":          pr,
		"replaced_by": replacedBy,
	})
}
