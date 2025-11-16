package http

import (
	"encoding/json"
	"net/http"
	"pull_request_service/internal/domain"
	"pull_request_service/internal/models"
)

type addTeamRequest struct {
	TeamName string `json:"team_name"`
	Members  []models.TeamMember
}

func (h *HTTPHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req addTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "invalid JSON")
		return
	}

	team, err := h.Team.CreateTeam(r.Context(), req.TeamName, req.Members)
	if err != nil {
		switch err {
		case domain.ErrTeamExists:
			writeError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"team": team,
	})
}

func (h *HTTPHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, http.StatusBadRequest, "BAD_INPUT", "team_name required")
		return
	}

	team, err := h.Team.GetTeam(r.Context(), teamName)
	if err != nil {
		if err == domain.ErrNotFound {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "team not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, team)
}
