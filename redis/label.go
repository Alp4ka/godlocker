package redis

import (
	"errors"
	"fmt"
	"github.com/Alp4ka/godlocker"
	"strings"
)

const (
	basePrefix  = "default"
	lockLabel   = "lock"
	sep         = "-"
	labelFormat = "%s" + sep + "%s" + sep + "%s"
)

const (
	prefixIdx = iota
	labelIdx
	idIdx
)

type Label string

func (l Label) Decompose() (prefix, id string, err error) {
	result := strings.Split(string(l), sep)
	ok, err := l.Valid()
	if err != nil || !ok {
		return "", "", errors.Join(godlocker.ErrDecomposeFailed, err)
	}

	return result[prefixIdx], result[idIdx], nil
}

func (l Label) String() string {
	return string(l)
}

func (l Label) Valid() (bool, error) {
	result := strings.Split(string(l), sep)

	if len(result) != 3 {
		return false, errors.Join(godlocker.ErrNotALabel, fmt.Errorf("expected 3 segments"))
	}

	if result[labelIdx] != lockLabel {
		return false, errors.Join(godlocker.ErrNotALabel, fmt.Errorf("expected '%s' as lock label, got: '%s'", lockLabel, result[labelIdx]))
	}

	return true, nil
}

func CreateLabel(prefix string, id string) (Label, error) {
	const (
		emptyString = ""
	)

	if strings.Compare(id, emptyString) == 0 || strings.Compare(strings.TrimSpace(id), emptyString) == 0 {
		return "", errors.Join(godlocker.ErrWrongLockerID, fmt.Errorf("expected not empty ID for locker label, got: '%s'", id))
	}

	if strings.Compare(prefix, emptyString) == 0 || strings.Compare(strings.TrimSpace(prefix), emptyString) == 0 {
		prefix = basePrefix
	}

	return Label(fmt.Sprintf(labelFormat, prefix, lockLabel, id)), nil
}

var _ godlocker.Label = (*Label)(nil)
