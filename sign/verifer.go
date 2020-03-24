package sign

import (
	"errors"
	"fmt"
	"github.com/guanaitong/crab/util/format"
	"time"
)

type Verifier struct {
	params     map[string]string
	secretKey  string
	signature  string
	cryptoFunc CryptoFunc
}

func NewVerifier(cryptoFunc CryptoFunc) *Verifier {
	return &Verifier{
		params:     map[string]string{},
		secretKey:  "",
		cryptoFunc: cryptoFunc,
	}
}

func NewVerifierDefault() *Verifier {
	return NewVerifier(HmacSign)
}

func (slf *Verifier) SetCryptoFunc(cryptoFunc CryptoFunc) *Verifier {
	slf.cryptoFunc = cryptoFunc
	return slf
}

func (slf *Verifier) SetParams(params map[string]string) *Verifier {
	slf.params = params
	return slf
}

func (slf *Verifier) SetQueryParams(args map[string]string) *Verifier {
	for k, v := range args {
		slf.params[k] = v
	}
	return slf
}

func (slf *Verifier) SetFormData(form map[string]string) *Verifier {
	for k, v := range form {
		slf.params[k] = v
	}
	return slf
}

func (slf *Verifier) SetBody(body interface{}) *Verifier {
	if body != nil {
		slf.params["x-req-body"] = format.AsString(body)
	}
	return slf
}

func (slf *Verifier) SetSecretKey(secretKey string) *Verifier {
	slf.secretKey = secretKey
	return slf
}

func (slf *Verifier) MustString(key string) string {
	if value, ok := slf.params[key]; ok {
		return value
	}
	return ""
}

func (slf *Verifier) MustInt64(key string) int64 {
	return convertToInt64(slf.MustString(key))
}

func (slf *Verifier) MustHasKeys(keys ...string) error {
	for _, key := range keys {
		if _, ok := slf.params[key]; !ok {
			return errors.New(fmt.Sprintf("KEY_MISSED:<%s>", key))
		}
	}
	return nil
}

func (slf *Verifier) CheckTimeStamp(param string, timeout time.Duration) error {
	timestamp := slf.MustInt64(param)
	thatTime := time.Unix(timestamp, 0)
	if time.Now().Sub(thatTime) > timeout {
		return errors.New(fmt.Sprintf("TIMESTAMP_TIMEOUT:<%d>", timestamp))
	}
	return nil
}

func (slf *Verifier) Check(signature string) bool {
	signer := NewSigner(slf.cryptoFunc).
		SetSecretKey(slf.secretKey).
		SetParams(slf.params)
	return signer.GetSignature() == signature
}
