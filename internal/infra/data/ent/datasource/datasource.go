package datasource

import (
	"context"

	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
)

type DataSource struct {
	client *model.Client
}

var _ tx.DataSource = (*DataSource)(nil)

func New(client *model.Client) *DataSource {
	return &DataSource{client: client}
}

func (ds *DataSource) Migrate(ctx context.Context) error {
	return ds.client.Schema.Create(ctx)
}

func (ds *DataSource) Key() any {
	return _txKey{}
}

func (ds *DataSource) NewTxFunc() tx.NewTxFunc {
	return func(ctx context.Context, opts tx.AtomicOpts) (tx.Tx, error) {
		return newTx(ctx, ds.client, opts)
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
