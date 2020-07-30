package gconf

import (
	"log"
)

// 该方法会在gconf后台同步goroutine里执行，请保证该方法不要有阻塞。不然会影响gconf更新。
// key      键
// oldValue 老的值,新增key时，该值为""
// newValue 新的值,删除key时，该值为""
type ConfigChangeListener interface {
	valueChanged(key, oldValue, newValue string)
}

// 配置集合
type ConfigCollection struct {
	appId     string
	name      string
	data      map[string]*Value //这里用map不线程安全不要紧，数据不会从map中移除，value指针会替换
	listeners map[string][]ConfigChangeListener
}

// 获取key对应的配置
func (c *ConfigCollection) GetValue(key string) *Value {
	res, ok := c.data[key]
	if ok {
		return res
	}
	return nil
}

// 获取配置结合中所有的key-value，以map返回。
func (c *ConfigCollection) AsMap() map[string]string {
	res := make(map[string]string)
	data := c.data // copy to avoid pointer change
	for k, v := range data {
		res[k] = v.Raw()
	}
	return res
}

func (c *ConfigCollection) AddConfigChangeListener(key string, configChangeListener ConfigChangeListener) {
	v, ok := c.listeners[key]
	if !ok {
		v = make([]ConfigChangeListener, 0, 1)
	}
	v = append(v, configChangeListener)
	c.listeners[key] = v
}

func (c *ConfigCollection) refreshData(client *gConfHttpClient) {
	newDataMap := client.listConfigs(c.appId)
	if len(newDataMap) == 0 {
		return
	}
	dataMap := c.data
	for key, oldValue := range dataMap {
		newValue, ok := newDataMap[key]
		if ok {
			o := oldValue.value
			if oldValue.refresh(newValue) {
				c.fireValueChanged(key, o, newValue)
			}
		} else { //老的有，但新的没有，先不从缓存里删除，避免程序出错。
			c.fireValueChanged(key, oldValue.value, "")
		}
	}
	for key, newV := range newDataMap {
		_, ok := dataMap[key]
		if !ok {
			dataMap[key] = newValue(key, newV)
		}
	}
}

func (c *ConfigCollection) fireValueChanged(key, oldValue, newValue string) {
	log.Printf("valueChanged,appId %s,key %s,oldValue--------->:\n%s\n    newValue--------->:\n%s", c.appId, key, oldValue, newValue)
	if listeners, ok := c.listeners[key]; ok {
		for _, listener := range listeners {
			listener.valueChanged(key, oldValue, newValue)
		}
	}
	log.Printf("firedValueChanged,appId %s,key %s", c.appId, key)
}
