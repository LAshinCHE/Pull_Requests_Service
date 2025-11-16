package postgres

import (
	"context"
	"pull_request_service/internal/domain"
	"pull_request_service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{pool: pool}
}

func (r *TeamRepo) Exists(ctx context.Context, teamID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM teams WHERE id = $1 
        )
    `, teamID).Scan(&exists)

	return exists, err
}

func (r *TeamRepo) Create(ctx context.Context, team *models.Team) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`, team.TeamName).
		Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrNotFound
	}

	_, err = tx.Exec(ctx, `
    INSERT INTO teams (team_name) VALUES ($1)
    ON CONFLICT (team_name) DO NOTHING
  `, team.TeamName)
	if err != nil {
		return err
	}

	for _, m := range team.Members {
		_, err = tx.Exec(ctx, `
      INSERT INTO users (id, username, team_name, is_active)
      VALUES ($1, $2, $3, $4)
      ON CONFLICT (id) DO UPDATE SET
        username = EXCLUDED.username,
        team_name = EXCLUDED.team_name,
        is_active = EXCLUDED.is_active
    `, m.UserID, m.Username, team.TeamName, m.IsActive)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `
      INSERT INTO team_members (team_name, id, is_active)
      VALUES ($1, $2, $3)
      ON CONFLICT (team_name, id) DO UPDATE SET
        is_active = EXCLUDED.is_active
    `, team.TeamName, m.UserID, m.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TeamRepo) GetByName(ctx context.Context, teamName string) (*models.Team, error) {
	var team models.Team
	team.TeamName = teamName

	exists, err := r.Exists(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrNotFound
	}

	rows, err := r.pool.Query(ctx, `
        SELECT id
        FROM team_members
        WHERE team_name = $1
    `, teamName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tm models.TeamMember
		tm.Username = teamName
		if err := rows.Scan(&tm.UserID); err != nil {
			return nil, err
		}
		team.Members = append(team.Members, tm)
	}

	return &team, nil
}

func (r *TeamRepo) ReplaceMembers(
	ctx context.Context,
	teamName string,
	members []models.TeamMember,
) error {

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM teams WHERE team_name = $1
        )
    `, teamName).Scan(&exists)

	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrNotFound
	}

	_, err = tx.Exec(ctx, `
        DELETE FROM team_members
        WHERE team_name = $1
    `, teamName)
	if err != nil {
		return err
	}

	for _, m := range members {
		_, err = tx.Exec(ctx, `
            INSERT INTO team_members (team_name, id)
            VALUES ($1, $2)
        `, teamName, m.UserID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
