package tx

import (
	"context"
	"errors"
	"fmt"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type _atomicKey struct{}

type Atomic struct{ txs map[any]Tx }

type AtomicOpts struct {
	ReadOnly    bool
	DataSources []any
}

func NewAtomic(ctx context.Context, opts AtomicOpts) (context.Context, error) {
	atomic := &Atomic{txs: make(map[any]Tx)}

	for _, v := range opts.DataSources {
		ds, ok := v.(DataSource)
		if !ok {
			// Skip in case it doesn't implement atomic operation.
			continue
		}
		if _, ok := atomic.txs[ds.Key()]; ok {
			continue
		}

		tx, err := ds.NewTxFunc()(ctx, opts)
		if err != nil {
			return nil, errors.Join(err, atomic.rollback(ctx))
		}

		atomic.txs[ds.Key()] = tx
	}

	return context.WithValue(ctx, _atomicKey{}, atomic), nil
}

func AtomicFromContext(ctx context.Context) *Atomic {
	atomic, ok := ctx.Value(_atomicKey{}).(*Atomic)
	if !ok {
		return nil
	}
	return atomic
}

func (a *Atomic) Get(key any) Tx {
	return a.txs[key]
}

func (a *Atomic) commit(ctx context.Context) error {
	errs := make([]error, 0)
	for _, tx := range a.txs {
		if e := tx.Commit(ctx); e != nil {
			errs = append(errs, e)
		}
	}
	return errors.Join(errs...)
}

func (a *Atomic) rollback(ctx context.Context) error {
	errs := make([]error, 0)
	for _, tx := range a.txs {
		if e := tx.Rollback(ctx); e != nil {
			errs = append(errs, e)
		}
	}
	return errors.Join(errs...)
}

func Evaluate(ctx context.Context, err *error) {
	if v := recover(); v != nil {
		err := fmt.Errorf("panic during atomic: %v", v)

		evaluate(ctx, &err)
		panic(err)
	}
	evaluate(ctx, err)
}

func evaluate(ctx context.Context, err *error) {
	atomic := AtomicFromContext(ctx)
	if atomic == nil {
		return
	}

	if *err == nil {
		*err = atomic.commit(ctx)
		return
	}

	rollbackErr := atomic.rollback(ctx)
	if rollbackErr != nil {
		*err = fmt.Errorf("failed to rollback atomic: %s, original: %w", rollbackErr.Error(), *err)
	}
}
