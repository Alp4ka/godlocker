package godlocker

import (
	"fmt"
)

const (
	_errLabel = "[locker]"
)

var (
	ErrLabelCreationFailed = fmt.Errorf("%s: label creation failed", _errLabel)
	ErrNotALabel           = fmt.Errorf("%s: not a label", _errLabel)
	ErrWrongLabelID        = fmt.Errorf("%s: label wrong id value", _errLabel)
	ErrWrongLabelPrefix    = fmt.Errorf("%s: label wrong prefix value", _errLabel)
	ErrWrongLabelName      = fmt.Errorf("%s: label wrong name", _errLabel)
	ErrDecomposeFailed     = fmt.Errorf("%s: label decomposition error", _errLabel)

	ErrReleaseFailed = fmt.Errorf("%s: mutex release fail", _errLabel)
	ErrAcquireFailed = fmt.Errorf("%s: mutex acquire fail", _errLabel)
)
