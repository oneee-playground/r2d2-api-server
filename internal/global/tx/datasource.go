package tx

import "context"

type NewTxFunc func(context.Context, AtomicOpts) (Tx, error)

type DataSource interface {
	Key() any
	NewTxFunc() NewTxFunc
}
