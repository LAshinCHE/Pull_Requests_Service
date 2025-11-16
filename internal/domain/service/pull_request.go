package service

import (
	"context"
	"pull_request_service/internal/domain"
	"pull_request_service/internal/models"
	"pull_request_service/internal/repository"
	"time"
)

type PullRequest struct {
	prRepo   repository.PullRequest
	userRepo repository.User
}

func NewPullRequest(pr repository.PullRequest, user repository.User) *PullRequest {
	return &PullRequest{
		prRepo:   pr,
		userRepo: user,
	}
}

func (p *PullRequest) GetPRsByReviewer(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	return p.prRepo.GetByReviewer(ctx, userID)
}

func (p *PullRequest) Create(ctx context.Context, prID string, prName string, authorID string) (*models.PullRequest, error) {

	_, err := p.prRepo.Get(ctx, prID)
	if err == nil {
		return nil, domain.ErrPRExists
	}

	reviewers, err := p.prRepo.GetAuthorTeamMembers(ctx, authorID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if len(reviewers) == 0 {
		return nil, domain.ErrNotFound
	}

	assigned := reviewers
	if len(assigned) > 2 {
		assigned = assigned[:2]
	}

	pr := &models.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            "OPEN",
		AssignedReviewers: assigned,
	}

	err = p.prRepo.Create(ctx, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (p *PullRequest) Merge(ctx context.Context, prID string) (*models.PullRequest, error) {
	exist, err := p.prRepo.Exists(ctx, prID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, domain.ErrNotFound
	}
	mergedTime := time.Now()
	err = p.prRepo.SetMerged(ctx, prID, &mergedTime)
	if err != nil {
		return nil, err
	}

	pr, err := p.prRepo.Get(ctx, prID)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (p *PullRequest) Reassign(
	ctx context.Context,
	prID string,
	oldUserID string,
) (*models.PullRequest, string, error) {
	return p.prRepo.Reassign(ctx, prID, oldUserID)
}
