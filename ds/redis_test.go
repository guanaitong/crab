package ds_test

import (
	"github.com/guanaitong/crab/ds"
	"github.com/guanaitong/crab/system"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	key      = "key"
	value    = "value"
	expected = "value"
)

func TestGetRedisConfig(t *testing.T) {
	system.SetupAppName("for-test-java")
	d := ds.GetDefaultRedisConfig()
	//d := ds.GetRedisConfig("redis-config.json")
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
	d := &ds.RedisConfig{Type: 1, Sentinel: ds.SentinelConfig{Nodes: "10.101.11.126:26379,10.101.11.127:26379,10.101.11.128:26379", Master: "mymaster"}}
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
