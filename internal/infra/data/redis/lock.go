package redis

import (
	"context"

	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/pkg/errors"
	"github.com/redis/rueidis/rueidislock"
)

type RedisLocker struct {
	underlying rueidislock.Locker
}

var _ tx.Locker = (*RedisLocker)(nil)

func NewLocker(lock rueidislock.Locker) *RedisLocker {
	return &RedisLocker{underlying: lock}
}

func (l *RedisLocker) AcquireKey(ctx context.Context, key string) (context.Context, context.CancelFunc, error) {
	ctx, release, err := l.underlying.WithContext(ctx, key)
	if err != nil {
		// TODO: What happens if err == rueidislock.ErrLockerClosed?
		return nil, nil, errors.Wrap(err, "acquiring key")
	}

	return ctx, release, nil
}

func (l *RedisLocker) Acquire(ctx context.Context, toJoin ...string) (context.Context, context.CancelFunc, error) {
	return l.AcquireKey(ctx, buildKey(toJoin...))
}
