package domain

import (
	"context"
	"pull_request_service/internal/models"
)

type ServiceTeam interface {
	CreateTeam(ctx context.Context, teamName string, members []models.TeamMemberInput) (models.Team, error)
	GetTeam(ctx context.Context, teamName string) (models.Team, error)
	GetTeamByUserID(ctx context.Context, userID string) (models.Team, error)
}

type ServiceUser interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (models.User, error)
	GetUser(ctx context.Context, userID string) (models.User, error)
}

type ServicePullRequest interface {
	Create(ctx context.Context, prID string, prName string, authorID string) (models.PullRequest, error)
	Merge(ctx context.Context, prID string) (models.PullRequest, error)
	Get(ctx context.Context, prID string) (models.PullRequest, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]models.PullRequest, error)
	Reassign(ctx context.Context, prID string, oldReviewerID string) (models.PullRequest, string, error)
}
