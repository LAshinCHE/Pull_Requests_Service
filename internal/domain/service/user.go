package service

import (
	"context"
	"pull_request_service/internal/domain"
	"pull_request_service/internal/models"
	"pull_request_service/internal/repository"
)

type User struct {
	userRepo repository.User
}

func NewUser(repo repository.User) *User {
	return &User{
		userRepo: repo,
	}
}

func (u *User) SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	exitst := u.userRepo.Exists(ctx, userID)
	if !exitst {
		return nil, domain.ErrNotFound
	}

	err := u.userRepo.SetActive(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
