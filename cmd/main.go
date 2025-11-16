package main

import (
	"context"
	"log"
	"os/signal"
	"pull_request_service/cmd/internal"
	"pull_request_service/internal/api/http"
	"pull_request_service/internal/config"
	"pull_request_service/internal/domain/service"
	postgres "pull_request_service/internal/repository/postgre"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	config := config.LoadConfig()

	db, err := internal.NewPostgreSQL(config.DBURL)
	if err != nil {
		log.Fatalf("Could not initialize Database connection %s", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepo(db)
	pullRequestRepo := postgres.NewPullRequestRepo(db)
	teamRepo := postgres.NewTeamRepo(db)

	prService := service.NewPullRequest(pullRequestRepo, userRepo)
	userService := service.NewUser(userRepo)
	teamService := service.NewTeam(teamRepo)

	http.RunService(ctx, userService, teamService, prService)

}
