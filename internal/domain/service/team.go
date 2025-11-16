package service

import (
	"context"
	"pull_request_service/internal/domain"
	"pull_request_service/internal/models"
	"pull_request_service/internal/repository"
)

type Team struct {
	teamRepo repository.Team
}

func NewTeam(team repository.Team) *Team {
	return &Team{
		teamRepo: team,
	}
}

func (t *Team) CreateTeam(ctx context.Context, teamName string, members []models.TeamMember) (*models.Team, error) {
	team := models.NewTeam(teamName, members)
	err := t.teamRepo.Create(ctx, team)
	if err != nil {
		return nil, domain.ErrTeamExists
	}
	return team, nil
}

func (t *Team) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	team, err := t.teamRepo.GetByName(ctx, teamName)

	if err != nil {
		return nil, domain.ErrNotFound
	}

	return team, nil
}
