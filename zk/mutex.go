package zk

import (
	"context"
	"errors"
	"github.com/Alp4ka/godlocker"
	"github.com/go-zookeeper/zk"
)

type Mutex struct {
	zkl *zk.Lock
}

func (m *Mutex) LockContext(ctx context.Context) error {
	err := m.zkl.Lock()
	if err != nil {
		return err
	}

	return nil
}

func (m *Mutex) UnlockContext(ctx context.Context) (bool, error) {
	err := m.zkl.Unlock()
	if err == nil {
		return true, nil
	}

	return false, errors.Join(godlocker.ErrReleaseFailed, err)
}

func NewMutex(zkl *zk.Lock) *Mutex {
	return &Mutex{zkl}
}

var _ godlocker.Mutex = (*Mutex)(nil)
