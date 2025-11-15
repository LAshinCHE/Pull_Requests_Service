package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(driver *pgxpool.Pool) *Repository {
	return &Repository{
		driver,
	}
}
