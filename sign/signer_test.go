package sign_test

import (
	"github.com/guanaitong/crab/sign"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	secretKey = "ZjcxZDUwZTRlZjViOTU5NTFkY2U1NGNhMDZmNmZhMGYK"
	signature = "d5ed54786712d8362c4def66f3509f99"
)

var (
	params = map[string]string{
		"contentType": "application/javascript; charset=utf8",
		"path":        "/apiserver-service/task/pull",
		"body":        "{}",
	}
	args = map[string]interface{}{
		"accessKey": "NzRjMWY1MmZmMjI5MmY4YjQyODc4N2Q3NTY3ODA1MjkK",
		"timestamp": "1584501049",
	}
)

func TestSigner_GetSignature(t *testing.T) {
	signer := sign.NewSignerDefault()
	signer.
		SetSecretKey(secretKey).
		SetParams(params).
		SetQuery(args)

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
			SetQuery(args).
			GetSignature()
	}
}
