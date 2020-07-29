package gconf

import (
	"log"
)

// 该方法会在gconf后台同步goroutine里执行，请保证该方法不要有阻塞。不然会影响gconf更新。
// key      键
// oldValue 老的值,新增key时，该值为""
// newValue 新的值,删除key时，该值为""
type ConfigChangeListener func(key, oldValue, newValue string)

// 配置集合
type Config struct {
	appId     string
	name      string
	data      map[string]*Value //这里用map不线程安全不要紧，数据不会从map中移除，value指针会替换
	listeners map[string][]ConfigChangeListener
}

// 获取key对应的配置
func (c *Config) GetValue(key string) *Value {
	res, ok := c.data[key]
	if ok {
		return res
	}
	return nil
}

// 获取配置结合中所有的key-value，以map返回。
func (c *Config) AsMap() map[string]string {
	res := make(map[string]string)
	data := c.data // copy to avoid pointer change
	for k, v := range data {
		res[k] = v.Raw()
	}
	return res
}

func (c *Config) AddConfigChangeListener(key string, configChangeListener ConfigChangeListener) {
	v, ok := c.listeners[key]
	if !ok {
		v = make([]ConfigChangeListener, 0, 1)
	}
	v = append(v, configChangeListener)
	c.listeners[key] = v
}

func (c *Config) refreshData() {
	newDataMap, err := getAppConfigs(c.appId)
	if err != nil {
		log.Println(err.Error())
		return
	}
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
			c.fireValueChanged(key, "", newV)
		}
	}
}

func (c *Config) fireValueChanged(key, oldValue, newValue string) {
	log.Printf("valueChanged,configCollectionId %s,key %s,oldValue:\n%s,newValue:\n%s", c.appId, key, oldValue, newValue)
	if listeners, ok := c.listeners[key]; ok {
		for _, listener := range listeners {
			listener(key, oldValue, newValue)
		}
	}
	log.Printf("firedValueChanged,configCollectionId %s,key %s", c.appId, key)
}

func getAppConfigs(appId string) (map[string]string, error) {
	dataMap, err := httpGetMapResp(baseUrl + "/api/listConfigs?configAppId=" + appId)
	if err != nil {
		return nil, err
	}
	structuredDataMap := make(map[string]string)
	for k, v := range dataMap {
		structuredDataMap[k] = v
	}
	return structuredDataMap, nil
}
