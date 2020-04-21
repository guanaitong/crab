package cache

import (
	"github.com/guanaitong/crab/errors"
	"github.com/guanaitong/crab/json"
	"github.com/guanaitong/crab/util/task"
	"k8s.io/klog"
	"sync"
	"time"
)

// Loader必须返回一个对象的指针
type Loader func() (interface{}, errors.Error)

type Cache interface {
	Get(key string, v interface{}, loader Loader) errors.Error
	Invalidate(key string) bool
}

type LocalCache struct {
	Data       map[string]*localValue
	Expiration time.Duration
}

var once sync.Once

var localCaches []*LocalCache

func NewLocalCache(Expiration time.Duration) Cache {
	cache := &LocalCache{Expiration: Expiration, Data: map[string]*localValue{}}
	localCaches = append(localCaches, cache)
	once.Do(initCleanTaskForLocalCache)
	return cache
}

func initCleanTaskForLocalCache() {
	task.StartBackgroundTask("clear_local_cache", time.Second*60, func() {
		for _, localCache := range localCaches {
			t := time.Now().Unix()
			var toDeleteKeys []string
			for k, v := range localCache.Data {
				if v.expireIn == 0 || t < v.expireIn {
					continue
				} else {
					toDeleteKeys = append(toDeleteKeys, k)
				}
			}
			for _, k := range toDeleteKeys {
				delete(localCache.Data, k)
			}
		}
	})
}

type localValue struct {
	value    []byte
	expireIn int64
}

func (c *LocalCache) Get(key string, v interface{}, loader Loader) errors.Error {
	t := time.Now().Unix()
	d, ok := c.Data[key]
	if ok {
		if d.expireIn == 0 || t < d.expireIn {
			err := json.Unmarshal(d.value, v)
			if err == nil {
				return nil
			} else {
				klog.Warningf("key %s value %s 反序列化失败", key, string(d.value))
			}
		} else { //数据已经过期
			c.Invalidate(key)
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
		klog.Warningf("key %s value %v 序列化失败", key, value)
		return errors.NewParamError(0, "序列化失败:"+err.Error())
	} else {
		var expireIn int64
		if c.Expiration > 0 {
			expireIn = int64(c.Expiration/time.Millisecond) + t
		}
		c.Data[key] = &localValue{value: bs, expireIn: expireIn}
		err = json.Unmarshal(bs, v)
		if err != nil {
			klog.Warningf("key %s value %s 反序列化失败", key, string(bs))
			return errors.NewParamError(0, err.Error())
		}
	}
	return nil
}

func (c *LocalCache) Invalidate(key string) bool {
	time.Now().Unix()
	_, ok := c.Data[key]
	if ok {
		delete(c.Data, key)
		return true
	}
	return false
}
