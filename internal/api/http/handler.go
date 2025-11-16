package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"pull_request_service/internal/domain"
)

func RunService(
	ctx context.Context,
	user domain.ServiceUser,
	team domain.ServiceTeam,
	pullRequest domain.ServicePullRequest,
) {
	handler := &HTTPHandler{
		User:        user,
		Team:        team,
		PullRequest: pullRequest,
	}

	r := newRouter(handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server ListenAndServe: %s", err)
	}
}

func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	resp := errorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	writeJSON(w, status, resp)
}

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type HTTPHandler struct {
	User        domain.ServiceUser
	Team        domain.ServiceTeam
	PullRequest domain.ServicePullRequest
}
