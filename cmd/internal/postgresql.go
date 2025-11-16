package internal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgreSQL(dburl string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dburl)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
