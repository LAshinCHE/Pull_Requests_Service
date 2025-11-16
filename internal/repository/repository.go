package repository

import (
	"context"
	"pull_request_service/internal/models"
	"time"
)

type User interface {
	Exists(ctx context.Context, userID string) bool
	CreateOrUpdate(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID string) (*models.User, error)
	SetActive(ctx context.Context, userID string, isActive bool) error
	GetTeamName(ctx context.Context, userID string) (string, error)
}

type Team interface {
	Create(ctx context.Context, team *models.Team) error
	Exists(ctx context.Context, teamName string) (bool, error)
	GetByName(ctx context.Context, teamName string) (*models.Team, error)

	ReplaceMembers(ctx context.Context, teamName string, members []models.TeamMember) error
}

type PullRequest interface {
	Create(ctx context.Context, pr *models.PullRequest) error
	Exists(ctx context.Context, prID string) (bool, error)
	Get(ctx context.Context, prID string) (*models.PullRequest, error)
	SetMerged(ctx context.Context, prID string, mergedAt *time.Time) error

	// UpdateReviewers(ctx context.Context, prID string, reviewers []string) error

	Reassign(ctx context.Context, prID string, oldUserID string) (*models.PullRequest, string, error)

	// для /users/getReview
	GetByReviewer(ctx context.Context, reviewerID string) ([]models.PullRequestShort, error)
}
