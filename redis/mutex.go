package redis

import (
	"context"
	"errors"
	"github.com/Alp4ka/godlocker"
	"github.com/go-redsync/redsync/v4"
)

type Mutex struct {
	mu *redsync.Mutex
}

func (m *Mutex) LockContext(ctx context.Context) error {
	err := m.mu.LockContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mutex) UnlockContext(ctx context.Context) (bool, error) {
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

func NewMutex(mu *redsync.Mutex) *Mutex {
	return &Mutex{mu}
}

var _ godlocker.Mutex = (*Mutex)(nil)
