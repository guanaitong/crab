package redis_test

import (
	"github.com/guanaitong/crab/cache"
	"github.com/guanaitong/crab/redis"
	"github.com/guanaitong/crab/system"
	"github.com/guanaitong/crab/util"
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

func TestLocalCache(t *testing.T) {
	system.SetupAppName("for-test-java")

	var c = redis.NewRedisCache("test", redis.GetDefaultRedisConfig().NewClient())
	c.Set("123", util.StringToBytes("456"), 0)
	v, err := c.Get("123")
	if err == nil && util.BytesToString(v) == "456" {
		t.Log("sucess")
	} else {
		t.Fail()
	}
	v, err = c.Get("1234")
	if err != cache.ErrEntryNotFound {
		t.Fail()
	}

	c.Set("123", util.StringToBytes("789"), time.Second)
	time.Sleep(time.Second * 2)
	v, err = c.Get("123")
	if err != cache.ErrEntryNotFound {
		t.Fail()
	}
}
