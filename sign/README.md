# AK/SK 签名

AK: `access_key`
SK: `secret_key`

## Usage

```bash
go get -u github.com/guanaitong/crab/sign
```

### Signer

```go
	signer := sign.NewSignerDefault()
	signer.
		SetSecretKey(secretKey).     // 设置SecretKey
		SetParams(params).           // 设置签名约定参数
		SetQueryParams(queryParams)  // 设置补充参数

	t.Log(signer.GetSignString()     // 获取签名体
	t.Log(signer.GetSignature())     // 获取签名串
	assert.EqualValues(t, signature, signer.GetSignature())

	t.Logf("%-6s: %s", "MD5", signer.SetCryptoFunc(sign.Md5Sign).GetSignature())
	t.Logf("%-6s: %s", "Hmac", signer.SetCryptoFunc(sign.HmacSign).GetSignature())
	t.Logf("%-6s: %s", "Sha256", signer.SetCryptoFunc(sign.Sha256Sign).GetSignature())

```

### Verifer

```go
	verifier := sign.NewVerifierDefault()
	verifier.
		SetSecretKey(secretKey).     // 设置SecretKey
		SetParams(params).           // 设置签名约定参数
		SetQueryParams(queryParams)  // 设置补充参数

	t.Log(verifier.Check(signature)) // 验证签名合法性
	assert.EqualValues(t, true, verifier.Check(signature))

```