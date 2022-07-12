package random

import (
	"math/rand"
	"time"
)

const (
	idChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	lenChar = 62
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		num := rand.Intn(lenChar)
		bytes[i] = idChars[num]
	}
	return string(bytes)
}

func RandomNumber(min int64, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}
