package sign_test

import (
	"github.com/guanaitong/crab/sign"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	secretKey = "ZjcxZDUwZTRlZjViOTU5NTFkY2U1NGNhMDZmNmZhMGYK"
	signature = "9f9e694600d2037573cf2d816f8fcbd7"
)

var (
	params map[string]string
)

func init() {
	params = map[string]string{
		"accessKey":   "NzRjMWY1MmZmMjI5MmY4YjQyODc4N2Q3NTY3ODA1MjkK",
		"timestamp":   "1584501049",
		"contentType": "application/javascript; charset=utf8",
		"uri":         "/apiserver-service/task/pull",
		"queryName":   "queryValue",
		"body":        "{}",
	}
}

func TestSigner_GetSignature(t *testing.T) {
	signer := sign.NewSignerDefault()
	signer.
		SetSecretKey(secretKey).
		SetParams(params)

	t.Log(signer.GetSignString())
	t.Log(signer.GetSignature())
	assert.EqualValues(t, signature, signer.GetSignature())

	t.Log(signer.SetCryptoFunc(sign.Md5Sign).GetSignature())
	t.Log(signer.SetCryptoFunc(sign.Sha256Sign).GetSignature())
}

func BenchmarkNewSignerHmac(b *testing.B) {
	for i := 0; i < b.N; i++ {
		signer := sign.NewSignerDefault()
		signer.
			SetSecretKey(secretKey).
			SetParams(params).
			GetSignature()
	}
}
