package godlocker

import "context"

type Mutex interface {
	LockContext(ctx context.Context) error
	UnlockContext(ctx context.Context) (bool, error)
}
