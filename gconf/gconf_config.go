package gconf

import (
	"errors"
	"log"
)

// 该方法会在gconf后台同步goroutine里执行，请保证该方法不要有阻塞。不然会影响gconf更新。
// key      键
// oldValue 老的值,新增key时，该值为""
// newValue 新的值,删除key时，该值为""
type ConfigChangeListener func(key, oldValue, newValue string)

type ConfigCollection struct {
	appId     string
	name      string
	data      map[string]*configData //这里用map不线程安全不要紧，数据不会从map中移除，value指针会替换
	listeners map[string][]ConfigChangeListener
}

func (c *ConfigCollection) getConfigData(key string) *configData {
	res, ok := c.data[key]
	if ok {
		return res
	}
	return nil
}

func (c *ConfigCollection) GetConfig(key string) string {
	v := c.getConfigData(key)
	if v == nil {
		return ""
	}
	return v.raw
}

//
func (c *ConfigCollection) GetConfigAsStructuredMap(key string) map[string]*Field {
	v := c.getConfigData(key)
	if v == nil {
		return map[string]*Field{}
	}
	return v.structuredData
}

// 将value赋值，value为struct指针。用法类似json.Unmarshal
func (c *ConfigCollection) GetConfigAsBean(key string, value interface{}) error {
	v := c.getConfigData(key)
	if v == nil {
		return errors.New("the value of key[" + key + "] is null")
	}
	return v.unmarshal(value)
}

func (c *ConfigCollection) AsMap() map[string]string {
	res := make(map[string]string)
	data := c.data // copy to avoid pointer change
	for k, v := range data {
		res[k] = v.raw
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

func (c *ConfigCollection) refreshData() {
	newDataMap, err := getStructuredDataMap(c.appId)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if len(newDataMap) == 0 {
		return
	}
	oldDataMap := c.data
	if len(oldDataMap) == 0 {
		c.data = newDataMap
		return
	}
	if configDataEqual(newDataMap, oldDataMap) {
		return
	}
	finalDataMap := map[string]*configData{}
	for key, oldValue := range oldDataMap {
		newValue, ok := newDataMap[key]
		if ok {
			if newValue.raw == oldValue.raw { //没改变，还是用老的值
				finalDataMap[key] = oldValue
			} else {
				finalDataMap[key] = newValue
				c.fireValueChanged(key, oldValue.raw, newValue.raw)

			}
		} else { //老的有，但新的没有，先不从缓存里删除，避免程序出错。
			finalDataMap[key] = oldValue
			c.fireValueChanged(key, oldValue.raw, "")
		}

	}
	for key, newValue := range newDataMap {
		_, ok := oldDataMap[key]
		if !ok {
			finalDataMap[key] = newValue
			c.fireValueChanged(key, "", newValue.raw)
		}
	}
	c.data = finalDataMap
}

func (c *ConfigCollection) fireValueChanged(key, oldValue, newValue string) {
	log.Printf("valueChanged,configCollectionId %s,key %s,oldValue:\n%s,newValue:\n%s", c.appId, key, oldValue, newValue)
	if listeners, ok := c.listeners[key]; ok {
		for _, listener := range listeners {
			listener(key, oldValue, newValue)
		}
	}
	log.Printf("firedValueChanged,configCollectionId %s,key %s", c.appId, key)
}

func configDataEqual(v1, v2 map[string]*configData) bool {
	if len(v1) != len(v2) {
		return false
	}
	for key, value := range v1 {
		vv, ok := v2[key]
		if !ok {
			return false
		}
		if vv.raw != value.raw {
			return false
		}
	}
	for key, value := range v2 {
		vv, ok := v1[key]
		if !ok {
			return false
		}
		if vv.raw != value.raw {
			return false
		}
	}
	return true
}

func getStructuredDataMap(appId string) (map[string]*configData, error) {
	dataMap, err := httpGetMapResp(baseUrl + "/api/listConfigs?configAppId=" + appId)
	if err != nil {
		return nil, err
	}
	structuredDataMap := make(map[string]*configData)
	for k, v := range dataMap {
		structuredDataMap[k] = toStructuredData(k, v)
	}
	return structuredDataMap, nil
}
