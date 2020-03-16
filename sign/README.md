# AK/SK 签名

AK: `access_key`
SK: `secret_key`

## Usage

```bash
go get -u github.com/guanaitong/crab/sign
```

### Signer

```go
	secretKey = "ZjcxZDUwZTRlZjViOTU5NTFkY2U1NGNhMDZmNmZhMGYK"
	params = map[string]string{
		"accessKey":   "NzRjMWY1MmZmMjI5MmY4YjQyODc4N2Q3NTY3ODA1MjkK",
		"timestamp":   "1584501049",
		"contentType": "application/javascript; charset=utf8",
		"uri":         "/apiserver-service/task/pull",
		"queryName":   "queryValue",
		"body":        "{}",
	}
	signer := sign.NewSignerDefault() // crypto/hmac 算法签名
	signer.
		SetSecretKey(secretKey).   // 设置secretKey
		SetParams(params).         // 设置签名参数，按约定参数体
		GetSignature()             // 获取签名串

```

### Verifer

```go
	secretKey = "ZjcxZDUwZTRlZjViOTU5NTFkY2U1NGNhMDZmNmZhMGYK"
	params = map[string]string{
		"accessKey":   "NzRjMWY1MmZmMjI5MmY4YjQyODc4N2Q3NTY3ODA1MjkK",
		"timestamp":   "1584501049",
		"contentType": "application/javascript; charset=utf8",
		"uri":         "/apiserver-service/task/pull",
		"queryName":   "queryValue",
		"body":        "{}",
	}
	verifier := sign.NewVerifierDefault()
	verifier.
	    SetSecretKey(secretKey).          // 设置SecretKey
		SetParams(params).                // 设置签名参数，按约定参数体
		Check(signature)                  // 验证签名合法性

```