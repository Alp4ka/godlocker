package utils

import (
	"math/rand"
	"time"
)

func RandSleep(from, to time.Duration) {
	rnd := rand.Int63()%(int64(to-from)) + int64(from)
	time.Sleep(time.Duration(rnd))
}
