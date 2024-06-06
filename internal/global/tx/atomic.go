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

func NewAtomic(ctx context.Context) context.Context {
	atomic := &Atomic{txs: make(map[any]Tx)}
	return context.WithValue(ctx, _atomicKey{}, atomic)
}

func AtomicFromContext(ctx context.Context) *Atomic {
	atomic, ok := ctx.Value(_atomicKey{}).(*Atomic)
	if !ok {
		return nil
	}
	return atomic
}

func (a *Atomic) Get(key any) (Tx, bool) {
	tx, ok := a.txs[key]
	return tx, ok
}

func (a *Atomic) Set(key any, tx Tx) {
	a.txs[key] = tx
}

func (a *Atomic) GetOrNew(key any, newFunc func() Tx) Tx {
	tx, ok := a.Get(key)
	if !ok {
		tx = newFunc()
		a.Set(key, tx)
	}

	return tx
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
