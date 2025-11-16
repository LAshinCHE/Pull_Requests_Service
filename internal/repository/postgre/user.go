package postgres

import (
	"context"
	"errors"
	"pull_request_service/internal/domain"
	"pull_request_service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Exists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM users WHERE id = $1
        )
    `, userID).Scan(&exists)

	return exists, err
}

func (r *UserRepo) CreateOrUpdate(ctx context.Context, user *models.User) error {

	_, err := r.pool.Exec(ctx, `
        INSERT INTO users (id, username, team_name, is_active)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id)
        DO UPDATE SET
            username = EXCLUDED.username,
            team_name = EXCLUDED.team_name,
            is_active = EXCLUDED.is_active
    `,
		user.UserID,
		user.Username,
		user.TeamName,
		user.IsActive,
	)

	return err
}

func (r *UserRepo) GetByID(ctx context.Context, userID string) (*models.User, error) {
	var u models.User

	err := r.pool.QueryRow(ctx, `
        SELECT
            user_id,
            username,
            team_name,
            is_active
        FROM users
        WHERE user_id = $1
    `, userID).Scan(
		&u.UserID,
		&u.Username,
		&u.TeamName,
		&u.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) SetActive(ctx context.Context, userID string, isActive bool) error {
	cmd, err := r.pool.Exec(ctx, `
        UPDATE users
        SET is_active = $1
        WHERE user_id = $2
    `, isActive, userID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *UserRepo) GetTeamName(ctx context.Context, userID string) (string, error) {
	var team string

	err := r.pool.QueryRow(ctx, `
        SELECT team_name
        FROM users
        WHERE user_id = $1
    `, userID).Scan(&team)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrNotFound
		}
		return "", err
	}

	return team, nil
}
