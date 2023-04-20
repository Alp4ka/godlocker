package redislocker

import "time"

// Busy wait.
const (
	// Lower constraint for busy wait.
	sleepLowerConstraint = time.Millisecond * 5

	// Upper constraint for busy wait.
	sleepUpperConstraint = time.Millisecond * 50
)

// Locker base settings.
const (
	// Default TTL for locker key in redis.
	defaultExpiry = time.Second * 5
)

// Label settings
const (
	_labelDefaultPrefix = "default"
	_labelLockRoot      = "lock"
	_labelNameSep       = "-"
	_labelNameFormat    = "%s" + _labelNameSep + "%s" + _labelNameSep + "%s"
)
