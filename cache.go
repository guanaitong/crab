package crab

import (
	"fmt"
	"sync"
	"time"
)

var globalCache = &container{
	caches:  make(map[string]*serviceCache),
	mux:     new(sync.Mutex),
	initial: false,
}

type container struct {
	caches  map[string]*serviceCache
	mux     *sync.Mutex
	initial bool
}

func (c *container) newServiceCache(serviceId string, client DiscoveryClient) (cache *serviceCache) {
	cache, ok := c.caches[serviceId]
	if ok {
		return
	}
	c.mux.Lock()
	defer c.mux.Unlock()

	cache, ok = c.caches[serviceId]
	if ok {
		return
	}
	cache = &serviceCache{
		serviceId:       serviceId,
		discoveryClient: client,
		instances:       make(map[string]*ServiceInstance),
		mux:             new(sync.Mutex),
	}
	cache.reload()
	c.caches[cache.serviceId] = cache

	if len(c.caches) == 1 {
		go func() {
			for {
				c.refresh()
				time.Sleep(time.Minute)
			}
		}()
	}
	return
}

func (c *container) refresh() {
	for _, v := range c.caches {
		v.reload()
	}
}

type serviceCache struct {
	serviceId       string
	discoveryClient DiscoveryClient
	instances       map[string]*ServiceInstance
	mux             *sync.Mutex
}

func (c *serviceCache) reloadAsync() {
	go c.reload()
}

func (c *serviceCache) reload() {
	newInstances, err := c.discoveryClient.GetInstances(c.serviceId)
	if err != nil {
		fmt.Printf("get newInstances error %s \n", err.Error())
		return
	}
	c.mux.Lock()
	defer c.mux.Unlock()

	// 最新获取的服务实例里没有了，把老的里删除
	for k := range c.instances {
		var existedInNew = false
		for _, newInstance := range newInstances {
			if k == newInstance.InstanceId {
				existedInNew = true
				break
			}
		}
		if !existedInNew {
			delete(c.instances, k)
		}
	}

	// add new
	for _, newInstance := range newInstances {
		_, ok := c.instances[newInstance.InstanceId]
		if !ok {
			c.instances[newInstance.InstanceId] = newInstance
		}
	}
}

func (c *serviceCache) GetUpdatedInstances() []*ServiceInstance {
	var res []*ServiceInstance
	for _, v := range c.instances {
		if !v.Status.NetFailed {
			res = append(res, v)
		}
	}
	//if res is empty,reload immediately
	if len(res) == 0 {
		c.reload()
		for _, v := range c.instances {
			if !v.Status.NetFailed {
				res = append(res, v)
			}
		}
	}
	return res
}
