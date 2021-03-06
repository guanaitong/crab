package cache_test

import (
	"github.com/guanaitong/crab/cache"
	"github.com/guanaitong/crab/util/strings2"
	"testing"
	"time"
)

func TestLocalCache(t *testing.T) {
	var c = cache.NewLocalCache()
	c.Set("123", strings2.StringToBytes("456"), 0)
	v, err := c.Get("123")
	if err == nil && strings2.BytesToString(v) == "456" {
		t.Log("sucess")
	} else {
		t.Fail()
	}
	v, err = c.Get("1234")
	if err != cache.ErrEntryNotFound {
		t.Fail()
	}

	c.Set("123", strings2.StringToBytes("789"), time.Second)
	time.Sleep(time.Second * 2)
	v, err = c.Get("123")
	if err != cache.ErrEntryNotFound {
		t.Fail()
	}
}
