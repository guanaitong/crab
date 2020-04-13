package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/guanaitong/crab/cache"
	"github.com/guanaitong/crab/json"
	"k8s.io/klog"
	"time"
)

type RedisCache struct {
	Client     *redis.Client
	Prefix     string
	Expiration time.Duration
}

func (c *RedisCache) Get(key string, v interface{}, loader cache.Loader) bool {
	redisKey := c.Prefix + key
	redisValue, err := c.Client.Get(redisKey).Bytes()
	if err != redis.Nil {
		if err != nil {
			klog.Warningf("get value of key (%s) from redis error:%s", redisKey, err.Error())
			return false
		} else {
			err := json.Unmarshal(redisValue, v)
			if err == nil {
				return true
			} else {
				klog.Warningf("key %s value %s 反序列化失败", redisKey, string(redisValue))
			}
		}
	}
	value := loader()
	if value == nil {
		return false
	}
	bs, err := json.Marshal(value)
	if err != nil {
		klog.Warningf("key %s value %v 序列化失败", redisKey, value)
	} else {
		err := c.Client.Set(redisKey, bs, c.Expiration).Err()
		if err != nil {
			klog.Warningf("set value of key (%s) to redis error:%s", redisKey, err.Error())
		}
		err = json.Unmarshal(bs, v)
		if err != nil {
			klog.Warningf("key %s value %s 反序列化失败", redisKey, string(redisValue))
		}
	}
	return false
}

func (c *RedisCache) Invalidate(key string) bool {
	redisKey := c.Prefix + key
	r := c.Client.Del(redisKey)
	if r.Err() != nil {
		return false
	}
	return r.Val() > 0
}
