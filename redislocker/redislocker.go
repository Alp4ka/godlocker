package redislocker

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

// Locker implementation of godlocker.Locker interface using redlock algorithm
type Locker struct {
	rs     *redsync.Redsync
	client *redis.Client

	// expiry default key expiration time. We use expiry divided by 2 to push signal to renew ttl of a key.
	expiry time.Duration
}

// mutex creates new godlocker.Mutex implementation to use it inside godlocker.Locker
func (l *Locker) mutex(label godlocker.Label) *Mutex {
	return NewMutex(
		l.rs.NewMutex(
			label.String(),
			redsync.WithExpiry(l.expiry),
			redsync.WithGenValueFunc(label.Hash),
		),
		l.expiry/2,
	)
}

// Lock creates mutex with the specified label from params and tries to lock it.
// Whether the withRetry param was set up with true it tries until success. I made it so
// to make this custom lock to act like std mutex.
func (l *Locker) Lock(ctx context.Context, label godlocker.Label, withRetry bool) (godlocker.Mutex, error) {
	if ok, err := label.Valid(); err != nil || !ok {
		return nil, err
	}

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
			utils.RandSleep(sleepLowerConstraint, sleepUpperConstraint)
			continue
		}

		return nil, errors.Join(godlocker.ErrAcquireFailed, err)
	}
}

// TryLock calls Lock method with param withRetry set up to false in purpose to avoid retry loop and return error
// in the case we are not able to acquire lock.
func (l *Locker) TryLock(ctx context.Context, label godlocker.Label) (godlocker.Mutex, error) {
	mu, err := l.Lock(ctx, label, false)
	if err != nil {
		return nil, err
	}

	return mu, nil
}

// Unlock unlocks mutex.
func (l *Locker) Unlock(ctx context.Context, mu godlocker.Mutex) error {
	if ok, err := mu.UnlockContext(ctx); !ok || err != nil {
		return errors.Join(godlocker.ErrReleaseFailed, err)

	}

	return nil
}

// CreateLabel creates label for specified Locker type. In this case we use redislocker.Label
func (l *Locker) CreateLabel(prefix string, id string) (godlocker.Label, error) {
	return CreateLabel(prefix, id)
}

// NewRedisLocker returns a pointer to redislocker.Locker with applied options and specified redisClient
func NewRedisLocker(redisClient *redis.Client, options ...Option) *Locker {
	pool := goredis.NewPool(redisClient)
	rs := redsync.New(pool)
	res := Locker{
		rs:     rs,
		client: redisClient,
		expiry: defaultExpiry,
	}

	for _, opt := range options {
		opt.Apply(&res)
	}

	return &res
}

var _ godlocker.Locker = (*Locker)(nil)
