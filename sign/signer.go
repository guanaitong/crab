package sign

import (
	"fmt"
	"sort"
	"strings"
)

type Signer struct {
	params     map[string]string
	secretKey  string
	bodyPrefix string
	bodySuffix string
	splitChar  string
	cryptoFunc CryptoFunc
}

func NewSigner(cryptoFunc CryptoFunc) *Signer {
	return &Signer{
		params:     map[string]string{},
		secretKey:  "",
		bodyPrefix: "",
		bodySuffix: "",
		splitChar:  "",
		cryptoFunc: cryptoFunc,
	}
}

func NewSignerDefault() *Signer {
	return NewSigner(HmacSign)
}

func (slf *Signer) SetCryptoFunc(cryptoFunc CryptoFunc) *Signer {
	slf.cryptoFunc = cryptoFunc
	return slf
}

func (slf *Signer) SetParams(params map[string]string) *Signer {
	slf.params = params
	return slf
}

func (slf *Signer) SetSecretKey(secretKey string) *Signer {
	slf.secretKey = secretKey
	return slf
}

func (slf *Signer) SetSecretKeyWrapBody(secretKey string) *Signer {
	slf.SetSignBodyPrefix(secretKey)
	slf.SetSignBodySuffix(secretKey)
	return slf.SetSecretKey(secretKey)
}

func (slf *Signer) SetSignBodyPrefix(prefix string) *Signer {
	slf.bodyPrefix = prefix
	return slf
}

func (slf *Signer) SetSignBodySuffix(suffix string) *Signer {
	slf.bodySuffix = suffix
	return slf
}

func (slf *Signer) SetSplitChar(split string) *Signer {
	slf.splitChar = split
	return slf
}

func (slf *Signer) GetSignString() string {
	return slf.bodyPrefix + slf.splitChar + slf.sortedParams() + slf.splitChar + slf.bodySuffix
}

func (slf *Signer) GetSignature() string {
	return fmt.Sprintf("%x", slf.cryptoFunc(slf.secretKey, slf.GetSignString()))
}

func (slf *Signer) sortedParams() string {
	size := len(slf.params)
	if size == 0 {
		return ""
	}

	keys := make([]string, size)
	idx := 0
	for k := range slf.params {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)

	pairs := make([]string, size)
	for i, key := range keys {
		pairs[i] = key + "=" + slf.params[key]
	}

	return strings.Join(pairs, "&")
}
