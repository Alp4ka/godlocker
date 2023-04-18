package zk

import (
	"context"
	"errors"
	"github.com/Alp4ka/godlocker"
	"github.com/Alp4ka/godlocker/utils"
	"github.com/go-zookeeper/zk"
	"time"
)

const (
	defaultLockerTimeout = 1 * time.Minute
	defaultZkTimeOut     = 20 * time.Second
)

type Locker struct {
	conn *zk.Conn
	acls []zk.ACL
}

func (l *Locker) mutex(label godlocker.Label) *Mutex {
	return NewMutex(zk.NewLock(l.conn, label.String(), l.acls))
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

		if errors.Is(err, zk.ErrDeadlock) {
			utils.RandSleep(utils.SleepLowerConstraint, utils.SleepUpperConstraint)
			continue
		}

		return nil, errors.Join(godlocker.ErrAcquireFailed, err)
	}
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

func (l *Locker) TryLock(ctx context.Context, label godlocker.Label) (godlocker.Mutex, error) {
	mu, err := l.Lock(ctx, label, false)
	if err != nil {
		return nil, err
	}

	return mu, nil
}

func (l *Locker) CreateLabel(prefix string, id string) (godlocker.Label, error) {
	return CreateLabel(prefix, id)
}

func NewZkLocker(conn *zk.Conn, options ...Option) *Locker {
	res := Locker{
		conn: conn,
		acls: zk.WorldACL(zk.PermAll),
	}

	for _, opt := range options {
		opt.Apply(&res)
	}

	return &res
}

var _ godlocker.Locker = (*Locker)(nil)
