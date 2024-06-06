package stubs

import "context"

type StubLocker struct{}

func NewStubLocker() *StubLocker { return &StubLocker{} }

func (StubLocker) AcquireKey(ctx context.Context, key string) (context.Context, context.CancelFunc, error) {
	return ctx, func() {}, nil
}

func (StubLocker) Acquire(ctx context.Context, toJoin ...string) (context.Context, context.CancelFunc, error) {
	return ctx, func() {}, nil
}
