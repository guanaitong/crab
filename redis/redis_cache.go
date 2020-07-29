package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/guanaitong/crab/cache"
	"time"
)

func NewRedisCache(prefix string, client *redis.Client) cache.Cache {
	return &redisCache{prefix: prefix, client: client}
}

type redisCache struct {
	client *redis.Client
	prefix string
}

func (c *redisCache) Get(key string) ([]byte, error) {
	redisKey := c.prefix + key
	redisValue, err := c.client.Get(redisKey).Bytes()
	if err == redis.Nil {
		return nil, cache.ErrEntryNotFound
	}
	if err != nil {
		return nil, err
	}
	return redisValue, nil
}

func (c *redisCache) Set(key string, entry []byte, ex time.Duration) error {
	redisKey := c.prefix + key
	return c.client.Set(redisKey, entry, ex).Err()
}

func (c *redisCache) Delete(key string) error {
	redisKey := c.prefix + key
	return c.client.Del(redisKey).Err()
}

func (c *redisCache) Reset() error {
	panic("implement me")
}
