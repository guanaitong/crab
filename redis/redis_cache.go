package redis

import (
	"fmt"
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

func (c *RedisCache) Get(key string, v interface{}, loader cache.Loader) error {
	redisKey := c.Prefix + key
	redisValue, err := c.Client.Get(redisKey).Bytes()
	if err != redis.Nil {
		if err != nil {
			klog.Warningf("get value of key (%s) from redis error:%s", redisKey, err.Error())
			return fmt.Errorf("get value of key (%s) from redis error:%w", redisKey, err)
		} else {
			err := json.Unmarshal(redisValue, v)
			if err == nil {
				return nil
			} else {
				klog.Warningf("key %s value %s 反序列化失败", redisKey, string(redisValue))
			}
		}
	}
	value, e := loader()
	if e != nil {
		return e
	}
	if value == nil {
		return nil
	}
	bs, err := json.Marshal(value)
	if err != nil {
		klog.Warningf("key %s value %v 序列化失败", redisKey, value)
		return fmt.Errorf("key %s value %v 序列化失败,%w", redisKey, value, err)
	} else {
		err := c.Client.Set(redisKey, bs, c.Expiration).Err()
		if err != nil {
			klog.Warningf("set value of key (%s) to redis error:%s", redisKey, err.Error())
		}
		err = json.Unmarshal(bs, v)
		if err != nil {
			klog.Warningf("key %s value %s 反序列化失败", redisKey, string(redisValue))
			return fmt.Errorf("key %s value %v 序列化失败,%w", redisKey, string(redisValue), err)
		}
		return nil
	}
}

func (c *RedisCache) Invalidate(key string) bool {
	redisKey := c.Prefix + key
	r := c.Client.Del(redisKey)
	if r.Err() != nil {
		return false
	}
	return r.Val() > 0
}
