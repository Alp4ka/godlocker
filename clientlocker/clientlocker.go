package clientlocker

import (
	"context"
	"github.com/Alp4ka/godlocker"
	"strconv"
	"sync"
)

const (
	prefix = "client"
)

var (
	_globalMu sync.RWMutex
	_globalCL *ClientLocker
)

type ClientLocker struct {
	locker godlocker.Locker
}

// Lock acquires mutex for client. If it is not possible to enter critical section, tries to do so until success.
func (l *ClientLocker) Lock(ctx context.Context, clientID int) (godlocker.Mutex, error) {
	label, err := l.locker.CreateLabel(prefix, strconv.Itoa(clientID))
	if err != nil {
		return nil, err
	}

	return l.locker.Lock(ctx, label, true)
}

// Unlock releases client mutex.
func (l *ClientLocker) Unlock(ctx context.Context, mu godlocker.Mutex) error {
	return l.locker.Unlock(ctx, mu)
}

// TryLock tries to acquire lock for specified client. Propagates error from Locker instance in case of fail.
// Returns godlocker.Mutex implementation when locker succeed.
func (l *ClientLocker) TryLock(ctx context.Context, clientID int) (godlocker.Mutex, error) {
	label, err := l.locker.CreateLabel(prefix, strconv.Itoa(clientID))
	if err != nil {
		return nil, err
	}

	return l.locker.TryLock(ctx, label)
}

// NewClientLocker creates new instance of ClientLocker.
func NewClientLocker(locker godlocker.Locker) *ClientLocker {
	return &ClientLocker{
		locker,
	}
}

// ReplaceGlobals use specified locker instance as global ClientLocker.
func ReplaceGlobals(locker *ClientLocker) {
	_globalMu.Lock()
	_globalCL = locker
	_globalMu.Unlock()
}

// CL gets global instance of ClientLocker.
func CL() *ClientLocker {
	_globalMu.RLock()
	cl := _globalCL
	_globalMu.RUnlock()
	return cl
}
