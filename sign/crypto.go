package sign

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
)

type CryptoFunc func(secretKey string, args string) []byte

func Md5Sign(_, body string) []byte {
	m := md5.New()
	m.Write([]byte(body))
	return m.Sum(nil)
}

func HmacSign(secretKey, body string) []byte {
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(body))
	s := h.Sum(nil)

	m := md5.New()
	m.Write(s)
	return m.Sum(nil)
}

func Sha256Sign(secretKey, body string) []byte {
	h := sha256.New()
	h.Write([]byte([]byte(body)))
	s := h.Sum([]byte(secretKey))

	m := md5.New()
	m.Write(s)
	return m.Sum(nil)
}
