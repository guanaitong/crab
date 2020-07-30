package strings2

import (
	"math/rand"
	"time"
	"unsafe"
)

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var randomSrc = rand.NewSource(time.Now().UnixNano())

func RandString(num int) string {
	if num < 4 {
		num = 4
	}
	if num > 1024 {
		num = 1024
	}
	bytes := make([]byte, num)
	for i, cache, remain := num-1, randomSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			bytes[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(bytes)
}
