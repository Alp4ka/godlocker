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

func (l *ClientLocker) Lock(ctx context.Context, clientID int) (godlocker.Mutex, error) {
	label, err := l.locker.CreateLabel(prefix, strconv.Itoa(clientID))
	if err != nil {
		return nil, err
	}

	return l.locker.Lock(ctx, label, true)
}

func (l *ClientLocker) Unlock(ctx context.Context, mu godlocker.Mutex) error {
	return l.locker.Unlock(ctx, mu)
}

func (l *ClientLocker) TryLock(ctx context.Context, clientID int) (godlocker.Mutex, error) {
	label, err := l.locker.CreateLabel(prefix, strconv.Itoa(clientID))
	if err != nil {
		return nil, err
	}

	return l.locker.TryLock(ctx, label)
}

func NewClientLocker(locker godlocker.Locker) *ClientLocker {
	return &ClientLocker{
		locker,
	}
}

func ReplaceGlobals(locker *ClientLocker) {
	_globalMu.Lock()
	_globalCL = locker
	_globalMu.Unlock()
}

func CL() *ClientLocker {
	_globalMu.RLock()
	cl := _globalCL
	_globalMu.RUnlock()
	return cl
}
