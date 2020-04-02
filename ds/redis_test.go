package ds

import (
	"github.com/guanaitong/crab/system"
	"testing"
)

func TestGetDefaultRedisConfig(t *testing.T) {
	system.SetupAppName("for-test-java")
	d := GetDefaultRedisConfig()
	client := d.NewClient()
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		t.Error(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		t.Error(err)
	}

	if val != "value" {
		t.Error("not equal")
	}
}

func TestSentinel(t *testing.T) {
	d := &RedisConfig{Type: 1, Sentinel: SentinelConfig{nodes: "10.101.11.126:26379,10.101.11.127:26379,10.101.11.128:26379", master: "mymaster"}}
	client := d.NewClient()
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		t.Error(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		t.Error(err)
	}

	if val != "value" {
		t.Error("not equal")
	}
}
