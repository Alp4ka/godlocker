package redislocker

import (
	"context"
	"errors"
	"github.com/Alp4ka/godlocker"
	"github.com/go-redsync/redsync/v4"
	"time"
)

type Mutex struct {
	mu         *redsync.Mutex
	extendRate time.Duration
	done       chan bool
}

func (m *Mutex) LockContext(ctx context.Context) error {
	err := m.mu.LockContext(ctx)
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(m.extendRate)
		for {
			select {
			case <-m.done:
				ticker.Stop()
				return
			case <-ticker.C:
				if ok, _ := m.mu.ExtendContext(ctx); !ok {
					ticker.Stop()
					return
				}
			}
		}
	}()

	return nil
}

func (m *Mutex) UnlockContext(ctx context.Context) (bool, error) {
	defer func() {
		go func() {
			m.done <- true
		}()
	}()
	ok, err := m.mu.UnlockContext(ctx)
	if err == nil && ok {
		return ok, nil
	}

	var et *redsync.ErrTaken
	if !errors.As(err, &et) {
		return ok, err
	}

	return false, errors.Join(godlocker.ErrReleaseFailed, err)
}

func NewMutex(mu *redsync.Mutex, extendRate time.Duration) *Mutex {
	return &Mutex{mu, extendRate, make(chan bool)}
}

var _ godlocker.Mutex = (*Mutex)(nil)
