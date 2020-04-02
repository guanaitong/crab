package ds_test

import (
	"github.com/guanaitong/crab/ds"
	"github.com/guanaitong/crab/system"
	"testing"
)

func TestGetDefaultRedisConfig(t *testing.T) {
	system.SetupAppName("for-test-java")
	d := ds.GetDefaultRedisConfig()
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
	d := &ds.RedisConfig{Type: 1, Sentinel: ds.SentinelConfig{Nodes: "10.101.11.126:26379,10.101.11.127:26379,10.101.11.128:26379", Master: "mymaster"}}
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
