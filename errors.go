package godlocker

import (
	"fmt"
)

const (
	errLabel = "locker"
)

var (
	ErrNotALabel       = fmt.Errorf("%s: not a label", errLabel)
	ErrWrongLockerID   = fmt.Errorf("%s: label wrong id value", errLabel)
	ErrDecomposeFailed = fmt.Errorf("%s: label decomposition error", errLabel)

	ErrReleaseFailed = fmt.Errorf("%s: mutex release fail", errLabel)
	ErrAcquireFailed = fmt.Errorf("%s: mutex acquire fail", errLabel)
)
