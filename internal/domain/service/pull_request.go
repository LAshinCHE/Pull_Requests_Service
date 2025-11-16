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

func (p *PullRequest) Create(ctx context.Context, pr *models.PullRequest) error {
	exist, err := p.prRepo.Exists(ctx, pr.PullRequestID)
	if err != nil {
		return err
	}

	if !exist {
		return domain.ErrPRExists
	}

	exist = p.userRepo.Exists(ctx, pr.AuthorID)
	if !exist {
		return domain.ErrNotFound
	}

	for _, temate := range pr.AssignedReviewers {
		if exist = p.userRepo.Exists(ctx, temate); !exist {
			return domain.ErrNotFound
		}
	}
	return p.Create(ctx, pr)
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
