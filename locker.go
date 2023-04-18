package godlocker

import (
	"context"
)

type Locker interface {
	TryLock(ctx context.Context, label Label) (Mutex, error)
	Lock(ctx context.Context, label Label, withRetry bool) (Mutex, error)
	Unlock(ctx context.Context, mu Mutex) error
	CreateLabel(prefix string, id string) (Label, error)
}
