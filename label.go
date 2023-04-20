package godlocker

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/Alp4ka/godlocker/utils"
)

const (
	_saltLength = 32
)

type Hashable interface {
	Hash() (string, error)
	Salt() string
}

type Label interface {
	String() string
	Decompose() (prefix, root, id string, err error)
	Valid() (ok bool, err error)
	Hashable
}

type HashLabel struct {
	salt string
}

func CreateHashLabel() HashLabel {
	return HashLabel{salt: utils.RandSeq(_saltLength)}
}

func (hl HashLabel) Salt() string {
	return hl.salt
}

func (hl HashLabel) Hash() (string, error) {
	h := sha1.New()
	h.Write([]byte(hl.Salt()))
	return hex.EncodeToString(h.Sum(nil)), nil
}

var _ Hashable = (*HashLabel)(nil)
