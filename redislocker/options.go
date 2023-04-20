package redislocker

import "time"

// Option is used to configure Locker to avoid constructor overloads
type Option interface {
	Apply(*Locker)
}

type OptionFunc func(locker *Locker)

// Apply applies option for the locker.
func (f OptionFunc) Apply(locker *Locker) {
	f(locker)
}

// WithExpiry is an option which defines TTL for locker key.
// It uses defaultExpiry when you pass incorrect expiry value to the function.
func WithExpiry(expiry time.Duration) Option {
	if expiry <= 0 {
		expiry = defaultExpiry
	}

	return OptionFunc(func(locker *Locker) {
		locker.expiry = expiry
	})
}
