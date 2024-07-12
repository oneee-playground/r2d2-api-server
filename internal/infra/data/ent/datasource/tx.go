package datasource

import (
	"context"
	"database/sql"

	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
)

type EntTx struct {
	underlying *model.Tx
}

func (tx *EntTx) Commit(ctx context.Context) error {
	return tx.underlying.Commit()
}

func (tx *EntTx) Rollback(ctx context.Context) error {
	return tx.underlying.Rollback()
}

type _txKey struct{}

func newTx(ctx context.Context, client *model.Client, opts tx.AtomicOpts) (*EntTx, error) {
	tx, err := client.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  opts.ReadOnly,
	})
	if err != nil {
		return nil, err
	}

	return &EntTx{underlying: tx}, nil
}

