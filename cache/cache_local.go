package cache

import (
	"sync"
	"time"
)

type localCache struct {
	data sync.Map
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
	v, ok := c.data.Load(key)
	if ok {
		if v.(*value).isExpired() { //数据已经过期
			c.Delete(key)
		} else {
			return v.(*value).value, nil
		}
	}
	return nil, ErrEntryNotFound
}

func (c *localCache) Set(key string, entry []byte, ex time.Duration) error {
	var expireIn int64
	if ex > 0 {
		expireIn = time.Now().Add(ex).Unix()
	}
	c.data.Store(key, &value{value: entry, expireIn: expireIn})
	return nil
}

func (c *localCache) Delete(key string) error {
	c.data.Delete(key)
	return nil
}

func (c *localCache) Reset() error {
	// new一个新map，赋予localCache，让老的map被gc
	c.data = sync.Map{}
	return nil
}

func (c *localCache) clearExpiredData() {
	var toDeleteKeys []interface{}
	c.data.Range(func(k, v interface{}) bool {
		if v.(*value).isExpired() {
			toDeleteKeys = append(toDeleteKeys, k)
		}
		return true
	})
	for _, k := range toDeleteKeys {
		c.data.Delete(k)
	}
}

var once sync.Once

var localCaches []*localCache

func NewLocalCache() Cache {
	cache := &localCache{data: sync.Map{}}
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
