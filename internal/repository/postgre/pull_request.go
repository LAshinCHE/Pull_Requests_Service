package postgres

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"pull_request_service/internal/domain"
	"pull_request_service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepo struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepo(pool *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{pool: pool}
}

func (r *PullRequestRepo) Get(ctx context.Context, id string) (*models.PullRequest, error) {
	var pr models.PullRequest

	err := r.pool.QueryRow(ctx, `
        SELECT id, author_id, status, assigned_reviewers
        FROM pull_requests
        WHERE id = $1
    `, id).Scan(
		&pr.PullRequestID,
		&pr.AuthorID,
		&pr.Status,
		&pr.AssignedReviewers,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &pr, nil
}

func (r *PullRequestRepo) Create(ctx context.Context, pr *models.PullRequest) error {
	_, err := r.pool.Exec(ctx, `
        INSERT INTO pull_requests (id, name, author_id, status, assigned_reviewers)
        VALUES ($1, $2, $3, $4, $5)
    `,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
		pr.AssignedReviewers,
	)
	return err
}

func (r *PullRequestRepo) Exists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM pull_requests WHERE id = $1
        )
    `, prID).Scan(&exists)

	return exists, err
}

func (r *PullRequestRepo) Update(ctx context.Context, pr models.PullRequest) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE pull_requests
        SET status = $1,
            assigned_reviewers = $2
        WHERE id = $3
    `,
		pr.Status,
		pr.AssignedReviewers,
		pr.PullRequestID,
	)
	return err
}

func (r *PullRequestRepo) SetMerged(ctx context.Context, prID string, mergedAt *time.Time) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE pull_requests
        SET status = 'MERGED',
            merged_at = $1
        WHERE id = $2
    `, mergedAt, prID)

	return err
}

func (r *PullRequestRepo) Reassign(ctx context.Context, prID string, oldUserID string) (*models.PullRequest, string, error) {

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback(ctx)

	var pr models.PullRequest

	err = tx.QueryRow(ctx, `
        SELECT id, author_id, status, assigned_reviewers
        FROM pull_requests
        WHERE id = $1
        FOR UPDATE
    `, prID).Scan(
		&pr.PullRequestID,
		&pr.AuthorID,
		&pr.Status,
		&pr.AssignedReviewers,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", domain.ErrNotFound
		}
		return nil, "", err
	}

	if pr.Status == models.StatusMerged {
		return nil, "", domain.ErrPRMerged
	}

	reviewerIndex := -1
	for i, r := range pr.AssignedReviewers {
		if r == oldUserID {
			reviewerIndex = i
			break
		}
	}
	if reviewerIndex == -1 {
		return nil, "", domain.ErrNotAssigned
	}

	var teamName string
	err = tx.QueryRow(ctx, `
        SELECT team_name
        FROM users
        WHERE user_id = $1
        FOR UPDATE
    `, oldUserID).Scan(&teamName)

	if err != nil {
		return nil, "", err
	}

	rows, err := tx.Query(ctx, `
        SELECT user_id
        FROM users
        WHERE team_name = $1
          AND is_active = TRUE
          AND user_id <> $2
          AND user_id <> $3
    `, teamName, oldUserID, pr.AuthorID)

	if err != nil {
		return nil, "", err
	}

	var candidates []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		candidates = append(candidates, id)
	}

	if len(candidates) == 0 {
		return nil, "", domain.ErrNoCandidate
	}

	newReviewer := candidates[rand.Intn(len(candidates))]

	pr.AssignedReviewers[reviewerIndex] = newReviewer

	_, err = tx.Exec(ctx, `
        UPDATE pull_requests
        SET assigned_reviewers = $1
        WHERE id = $2
    `, pr.AssignedReviewers, pr.PullRequestID)

	if err != nil {
		return nil, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", err
	}

	return &pr, newReviewer, nil
}

func (r *PullRequestRepo) GetByReviewer(ctx context.Context, reviewerID string) ([]models.PullRequestShort, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT
            id,
            name,
            author_id,
            status
        FROM pull_requests
        WHERE $1 = ANY(assigned_reviewers)
        ORDER BY created_at DESC
    `, reviewerID)

	if err != nil {
		return nil, err
	}

	var list []models.PullRequestShort

	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
		); err != nil {
			return nil, err
		}

		list = append(list, pr)
	}

	return list, nil
}

func (r *PullRequestRepo) GetAuthorTeamMembers(ctx context.Context, authorID string) ([]string, error) {
	const q = `
        SELECT tm.user_id
        FROM user_teams ut
        JOIN team_members tm ON tm.team_id = ut.team_id
        WHERE ut.user_id = $1 AND tm.user_id != $1
    `
	rows, err := r.pool.Query(ctx, q, authorID)
	if err != nil {
		return nil, err
	}

	reviewers := []string{}
	for rows.Next() {
		var uid string
		rows.Scan(&uid)
		reviewers = append(reviewers, uid)
	}

	return reviewers, nil
}
