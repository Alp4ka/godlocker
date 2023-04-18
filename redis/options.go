package redis

import "time"

type Option interface {
	Apply(*Locker)
}

type OptionFunc func(locker *Locker)

func (f OptionFunc) Apply(locker *Locker) {
	f(locker)
}

func WithExpiry(expiry time.Duration) Option {
	return OptionFunc(func(locker *Locker) {
		locker.expiry = expiry
	})
}
