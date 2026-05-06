package order

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool      *pgxpool.Pool
	txManager trm.Manager
	getter    *trmpgx.CtxGetter
}

func New(pool *pgxpool.Pool, txManager trm.Manager) *repository {
	return &repository{
		pool:      pool,
		txManager: txManager,
		getter:    trmpgx.DefaultCtxGetter,
	}
}
