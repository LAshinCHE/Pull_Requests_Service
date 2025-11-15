package http

import "github.com/gorilla/mux"

func NewRouter(h *HTTPHandler) *mux.Router {
	r := mux.NewRouter()

	// Team
	r.HandleFunc("/team/add", h.AddTeam).Methods("POST")
	r.HandleFunc("/team/get", h.GetTeam).Methods("GET")

	// Users
	r.HandleFunc("/users/setIsActive", h.SetUserIsActive).Methods("POST")
	r.HandleFunc("/users/getReview", h.GetUserReviews).Methods("GET")

	// Pull Requests
	r.HandleFunc("/pullRequest/create", h.CreatePullRequest).Methods("POST")
	r.HandleFunc("/pullRequest/merge", h.MergePullRequest).Methods("POST")
	r.HandleFunc("/pullRequest/reassign", h.ReassignReviewer).Methods("POST")

	// Healthcheck
	r.HandleFunc("/health", h.Health).Methods("GET")

	return r
}
