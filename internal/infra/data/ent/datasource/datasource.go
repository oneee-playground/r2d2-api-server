package datasource

import (
	"context"

	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
)

type DataSource struct {
	client *model.Client
}

func (ds *DataSource) Key() any {
	return _txKey{}
}

func (ds *DataSource) NewTxFunc() tx.NewTxFunc {
	return func(ctx context.Context, opts tx.AtomicOpts) (tx.Tx, error) {
		return New(ctx, ds.client, opts)
	}
}

func (ds *DataSource) TxOrPlain(ctx context.Context) *model.Client {
	atomic := tx.AtomicFromContext(ctx)
	if atomic == nil {
		return ds.client
	}

	tx, ok := atomic.Get(_txKey{}).(*EntTx)
	if !ok {
		return ds.client
	}
	return tx.underlying.Client()
}
