package order

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}
