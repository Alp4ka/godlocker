package redislocker

import (
	"errors"
	"fmt"
	"github.com/Alp4ka/godlocker"
	"strings"
)

const (
	_prefixIdx = iota
	_rootIdx
	_idIdx
)

type Label struct {
	string
	godlocker.HashLabel
}

// Decompose is needed to get prefix, root and id values(separately) from existing Label.
func (l Label) Decompose() (string, string, string, error) {
	result := strings.Split(l.String(), _labelNameSep)

	if len(result) != 3 {
		return "", "", "", errors.Join(godlocker.ErrDecomposeFailed, fmt.Errorf("expected 3 segments"))
	}

	return result[_prefixIdx], result[_rootIdx], result[_idIdx], nil
}

// String gets string value of Label
func (l Label) String() string {
	return l.string
}

// Valid validates Label format
func (l Label) Valid() (bool, error) {
	const (
		empty = ""
	)

	prefix, root, id, err := l.Decompose()
	if err != nil {
		return false, errors.Join(godlocker.ErrNotALabel, err)
	}

	if root != _labelLockRoot {
		return false, errors.Join(godlocker.ErrNotALabel, fmt.Errorf("expected '%s' as lock label, got: '%s'", _labelLockRoot, root))
	}

	if strings.Compare(strings.TrimSpace(id), empty) == 0 {
		return false, errors.Join(godlocker.ErrWrongLabelID, fmt.Errorf("expected not empty id for locker label, got: '%s'", id))
	}

	if strings.Compare(strings.TrimSpace(prefix), empty) == 0 {
		return false, errors.Join(godlocker.ErrWrongLabelPrefix, fmt.Errorf("expected not empty prefix for locker label, got: '%s'", prefix))
	}

	if strings.Contains(prefix+id, _labelNameSep) {
		return false, errors.Join(godlocker.ErrWrongLabelName, fmt.Errorf("name should not contain substring: '%s'", _labelNameSep))
	}

	return true, nil
}

// Hash we use it to store random key value in redis.
func (l Label) Hash() (string, error) {
	return l.HashLabel.Hash()
}

func CreateLabel(prefix string, id string) (Label, error) {
	const (
		emptyString = ""
	)

	// Defaults.
	if strings.Compare(prefix, emptyString) == 0 || strings.Compare(strings.TrimSpace(prefix), emptyString) == 0 {
		prefix = _labelDefaultPrefix
	}

	// Create label
	label := Label{
		fmt.Sprintf(_labelNameFormat, prefix, _labelLockRoot, id),
		godlocker.CreateHashLabel(),
	}

	// Validation
	if ok, err := label.Valid(); !ok || err != nil {
		return Label{}, errors.Join(godlocker.ErrLabelCreationFailed, err)
	}

	return label, nil
}

var _ godlocker.Label = (*Label)(nil)
