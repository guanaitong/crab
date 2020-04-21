package redis_test

import (
	"fmt"
	cache2 "github.com/guanaitong/crab/cache"
	errors2 "github.com/guanaitong/crab/errors"
	"github.com/guanaitong/crab/redis"
	"github.com/guanaitong/crab/system"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	key      = "key"
	value    = "value"
	expected = "value"
)

func TestGetRedisConfig(t *testing.T) {
	system.SetupAppName("for-test-java")
	d := redis.GetDefaultRedisConfig()
	//d := mysql.GetRedisConfig("redis-config.json")
	assert.NotNil(t, d)

	client := d.NewClient()
	err := client.Set(key, value, 0).Err()
	if !assert.NoError(t, err) {
		t.Error(err)
		return
	}

	actual, err := client.Get("key").Result()
	if !assert.NoError(t, err) {
		t.Error(err)
		return
	}

	assert.EqualValues(t, expected, actual)
}

func TestSentinel(t *testing.T) {
	d := &redis.RedisConfig{Type: 1, Sentinel: redis.SentinelConfig{Nodes: "10.101.11.126:26379,10.101.11.127:26379,10.101.11.128:26379", Master: "mymaster"}}
	client := d.NewClient()
	err := client.Set(key, value, 0).Err()
	if !assert.NoError(t, err) {
		t.Error(err)
		return
	}

	actual, err := client.Get("key").Result()
	if !assert.NoError(t, err) {
		t.Error(err)
		return
	}

	assert.EqualValues(t, expected, actual)
}

type User struct {
	Id   int
	Name string
	Time time.Time
}

func TestRedisCache_Get(t *testing.T) {
	system.SetupAppName("for-test-java")
	var cache cache2.Cache = &redis.RedisCache{Client: redis.GetDefaultRedisConfig().NewClient(), Prefix: "Test", Expiration: time.Hour}
	cache.Invalidate("1")
	user := new(User)
	b := cache.Get("1", user, func() (interface{}, errors2.Error) {
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

	b = cache.Get("1", user2, func() (interface{}, errors2.Error) {
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
