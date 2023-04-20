package utils

import (
	cryptorand "crypto/rand"
	"encoding/base32"
	mathrand "math/rand"
	"time"
)

func RandSleep(from, to time.Duration) {
	rnd := mathrand.Int63()%(int64(to-from)) + int64(from)
	time.Sleep(time.Duration(rnd))
}

func RandSeq(length int) string {
	randomBytes := make([]byte, length)
	_, err := cryptorand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}
