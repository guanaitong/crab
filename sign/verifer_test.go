package sign_test

import (
	"github.com/guanaitong/crab/sign"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifier_Check(t *testing.T) {
	verifier := sign.NewVerifierDefault()
	verifier.
		SetSecretKey(secretKey).
		SetParams(params).
		SetArgs(args).
		SetForm(form).
		SetData(data)

	t.Log(verifier.Check(signature))
	assert.EqualValues(t, true, verifier.Check(signature))
}

func BenchmarkVerifier_Check(b *testing.B) {
	for i := 0; i < b.N; i++ {
		verifier := sign.NewVerifierDefault()
		verifier.
			SetSecretKey(secretKey).
			SetParams(params).
			SetArgs(args).
			SetForm(form).
			SetData(data)

		verifier.Check(signature)
	}
}
