package tx

import "context"

type Locker interface {
	// AcquireKey aquires lock from given key.
	AcquireKey(ctx context.Context, key string) (context.Context, context.CancelFunc, error)
	// Acquire is like AcquireKey. But it joins toJoin with seperator ':' and useses it as key.
	Acquire(ctx context.Context, toJoin ...string) (context.Context, context.CancelFunc, error)
}
