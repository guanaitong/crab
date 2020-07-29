package cache

import (
	"sync"
	"time"
)

type localCache struct {
	data map[string]*value
}

type value struct {
	value    []byte
	expireIn int64
}

func (v *value) isExpired() bool {
	if v.expireIn == 0 {
		return false
	}
	return time.Now().Unix() > v.expireIn
}

func (c *localCache) Get(key string) ([]byte, error) {
	v, ok := c.data[key]
	if ok {
		if v.isExpired() { //数据已经过期
			c.Delete(key)
		} else {
			return v.value, nil
		}
	}
	return nil, ErrEntryNotFound
}

func (c *localCache) Set(key string, entry []byte, ex time.Duration) error {
	var expireIn int64
	if ex > 0 {
		expireIn = time.Now().Add(ex).Unix()
	}
	c.data[key] = &value{value: entry, expireIn: expireIn}
	return nil
}

func (c *localCache) Delete(key string) error {
	delete(c.data, key)
	return nil
}

func (c *localCache) Reset() error {
	// new一个新map，赋予localCache，让老的map被gc
	c.data = map[string]*value{}
	return nil
}

func (c *localCache) clearExpiredData() {
	var toDeleteKeys []string
	for k, v := range c.data {
		if v.isExpired() {
			toDeleteKeys = append(toDeleteKeys, k)
		}
	}
	for _, k := range toDeleteKeys {
		delete(c.data, k)
	}
}

var once sync.Once

var localCaches []*localCache

func NewLocalCache() Cache {
	cache := &localCache{data: map[string]*value{}}
	localCaches = append(localCaches, cache)
	once.Do(initCleanTaskForLocalCache)
	return cache
}

func initCleanTaskForLocalCache() {
	go func() {
		for {
			for _, localCache := range localCaches {
				localCache.clearExpiredData()
			}
			time.Sleep(time.Minute)
		}
	}()
}
