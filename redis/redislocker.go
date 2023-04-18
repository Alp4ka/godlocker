package redis

import (
	"context"
	"errors"
	"github.com/Alp4ka/godlocker"
	"github.com/Alp4ka/godlocker/utils"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	defaultExpiry = time.Minute * 5
)

type Locker struct {
	rs     *redsync.Redsync
	expiry time.Duration
}

func (l *Locker) mutex(label godlocker.Label) *Mutex {
	var expOpt redsync.Option

	if l.expiry <= 0 {
		expOpt = redsync.WithExpiry(defaultExpiry)
	} else {
		expOpt = redsync.WithExpiry(l.expiry)
	}

	return NewMutex(l.rs.NewMutex(label.String(), expOpt))
}

func (l *Locker) Lock(ctx context.Context, label godlocker.Label, withRetry bool) (godlocker.Mutex, error) {
	mu := l.mutex(label)
	if !withRetry {
		err := mu.LockContext(ctx)

		if err != nil {
			return nil, errors.Join(godlocker.ErrAcquireFailed, err)
		}
		return mu, nil
	}

	for {
		err := mu.LockContext(ctx)
		if err == nil {
			return mu, nil
		}

		var et *redsync.ErrTaken
		if errors.As(err, &et) {
			utils.RandSleep(utils.SleepLowerConstraint, utils.SleepUpperConstraint)
			continue
		}

		return nil, errors.Join(godlocker.ErrAcquireFailed, err)
	}
}

func (l *Locker) TryLock(ctx context.Context, label godlocker.Label) (godlocker.Mutex, error) {
	mu, err := l.Lock(ctx, label, false)
	if err != nil {
		return nil, err
	}

	return mu, nil
}

func (l *Locker) Unlock(ctx context.Context, mu godlocker.Mutex) error {
	var retErr error
	if ok, err := mu.UnlockContext(ctx); !ok {
		retErr = godlocker.ErrReleaseFailed
		if err != nil {
			retErr = errors.Join(retErr, err)
		}
	}

	if retErr != nil {
		return retErr
	}
	return nil
}

func (l *Locker) CreateLabel(prefix string, id string) (godlocker.Label, error) {
	return CreateLabel(prefix, id)
}

func NewRedisLocker(redisClient *redis.Client, options ...Option) *Locker {
	pool := goredis.NewPool(redisClient)
	rs := redsync.New(pool)
	res := Locker{
		rs: rs,
	}

	for _, opt := range options {
		opt.Apply(&res)
	}

	return &res
}

var _ godlocker.Locker = (*Locker)(nil)
