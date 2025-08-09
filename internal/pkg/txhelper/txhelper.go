package txhelper

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxHelper struct {
	pool *pgxpool.Pool
}

func NewTxHelper(pool *pgxpool.Pool) *TxHelper {
	return &TxHelper{pool}
}

func (txh *TxHelper) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := txh.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (txh *TxHelper) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	tx, err := txh.pool.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
