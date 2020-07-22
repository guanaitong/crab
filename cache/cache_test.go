package cache_test

import (
	"fmt"
	cache2 "github.com/guanaitong/crab/cache"
	"testing"
	"time"
)

type User struct {
	Id   int
	Name string
	Time time.Time
}

func TestRedisCache_Get(t *testing.T) {
	var cache = cache2.NewLocalCache(time.Minute)
	cache2.NewLocalCache(time.Minute)

	cache.Invalidate("1")
	user := new(User)
	b := cache.Get("1", user, func() (interface{}, error) {
		return &User{
			Id:   123456789,
			Name: "august",
			Time: time.Now(),
		}, nil
	})
	if b != nil {
		t.Fail()
	}
	user2 := new(User)

	b = cache.Get("1", user2, func() (interface{}, error) {
		return &User{
			Id:   123456789,
			Name: "august",
			Time: time.Now(),
		}, nil
	})
	if b != nil {
		t.Fail()
	}

	if *user != *user2 {
		t.Fail()
	}
	fmt.Println(b)
}
