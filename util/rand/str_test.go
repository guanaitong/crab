package rand_test

import (
	"github.com/guanaitong/crab/util/rand"
	"testing"
)

func TestRandString(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Log(rand.RandString(5))
	}
}

func BenchmarkRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.RandString(5)
	}
}
