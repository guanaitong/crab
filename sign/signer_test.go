package sign_test

import (
	"fmt"
	"github.com/guanaitong/crab/json"
	"github.com/guanaitong/crab/sign"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	accessKey = "NzRjMWY1MmZmMjI5MmY4YjQyODc4N2Q3NTY3ODA1MjkK"
	secretKey = "ZjcxZDUwZTRlZjViOTU5NTFkY2U1NGNhMDZmNmZhMGYK"
	signature = "a5349d3e067ba50a82a8d307edfdb514"
)

var (
	params = map[string]string{
		"x-sec-timestamp":    "1584501049",
		"x-sec-access-key":   accessKey,
		"x-req-content-type": "test-plain",
		"x-req-path":         "/apiserver-service/isLive",
		"x-req-method":       "GET",
		"x-req-length":       "0",
	}
	queryParams = map[string]string{
		"accessKey": "NzRjMWY1MmZmMjI5MmY4YjQyODc4N2Q3NTY3ODA1MjkK",
		"timestamp": "1584501049",
	}
	formData = map[string]string{
		"accessKey": "NzRjMWY1MmZmMjI5MmY4YjQyODc4N2Q3NTY3ODA1MjkK",
		"timestamp": "1584501049",
	}
	body = map[string]interface{}{
		"fruits":   []string{"apple", "pear", "banana", "mango"},
		"animals":  []string{"tiger", "lion"},
		"location": "Shanghai",
	}
)

func TestSigner_GetSignature(t *testing.T) {
	signer := sign.NewSignerDefault()
	signer.
		SetSecretKey(secretKey).
		SetParams(params).
		SetQueryParams(queryParams)

	t.Log(signer.GetSignString())
	t.Log(signer.GetSignature())
	assert.EqualValues(t, signature, signer.GetSignature())

	t.Logf("%-6s: %s", "MD5", signer.SetCryptoFunc(sign.Md5Sign).GetSignature())
	t.Logf("%-6s: %s", "Hmac", signer.SetCryptoFunc(sign.HmacSign).GetSignature())
	t.Logf("%-6s: %s", "Sha256", signer.SetCryptoFunc(sign.Sha256Sign).GetSignature())
}

func BenchmarkNewSignerHmac(b *testing.B) {
	for i := 0; i < b.N; i++ {
		signer := sign.NewSignerDefault()
		signer.
			SetSecretKey(secretKey).
			SetParams(params).
			SetQueryParams(queryParams)
	}
}

func TestSignerDiffMethod(t *testing.T) {
	//GET
	params["x-req-content-type"] = "test-plain"
	params["x-req-method"] = "GET"
	params["x-req-length"] = "0"
	t.Logf("%-4s: %s", params["x-req-method"],
		sign.
			NewSignerDefault().
			SetSecretKey(secretKey).
			SetParams(params).
			SetQueryParams(queryParams).
			GetSignature())

	//POST: form
	params["x-req-content-type"] = "application/x-www-form-urlencoded"
	params["x-req-method"] = "POST"
	params["x-req-length"] = fmt.Sprint(getFormDataLength(formData))
	t.Logf("%-4s: %s", params["x-req-method"],
		sign.
			NewSignerDefault().
			SetSecretKey(secretKey).
			SetParams(params).
			GetSignature())

	//POST: json
	params["x-req-content-type"] = "application/json; charset=UTF-8"
	params["x-req-method"] = "POST"
	params["x-req-length"] = fmt.Sprint(len(json.AsString(body)))
	t.Logf("%-4s: %s", params["x-req-method"],
		sign.
			NewSignerDefault().
			SetSecretKey(secretKey).
			SetParams(params).
			GetSignature())
}

func getFormDataLength(form map[string]string) int {
	length := -1
	for k, v := range form {
		length += len(k+v) + 2
	}
	return length
}
