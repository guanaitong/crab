package gconf

import (
	"github.com/guanaitong/crab/system"
	"github.com/guanaitong/crab/util/task"
	"log"
	"strings"
	"sync"
	"time"
)

var cache = map[string]*ConfigCollection{}
var baseUrl string
var mux = new(sync.Mutex)

func init() {
	if system.InK8s() {
		baseUrl = "http://gconf.kube-system"
	} else {
		baseUrl = "http://gconf" + system.GetServiceDomainSuffix()
	}

	task.StartBackgroundTask("gconf-refresher", time.Millisecond*100, func() {
		if len(cache) == 0 {
			time.Sleep(time.Second * 2)
			return
		}
		var keys []string
		for k := range cache {
			keys = append(keys, k)
		}
		configAppIdList := strings.Join(keys, ",")
		needChangeAppIdList, err := httpGetListResp(baseUrl + "/api/watch?configAppIdList=" + configAppIdList + "&clientId=" + system.GetInstanceId())
		if err != nil {
			log.Printf("wath error" + err.Error())
			time.Sleep(time.Second * 10)
			return
		}

		for _, appId := range needChangeAppIdList {
			cache[appId].refreshData()
		}
	})
}

// 获取当前应用的配置集合
func GetCurrentConfigCollection() *ConfigCollection {
	return GetConfigCollection(system.GetAppName())
}

// 获取全局的配置配置集合，此方法用于框架的统一配置。
// 应用不需要调用此方法
func GetGlobalConfigCollection() *ConfigCollection {
	return GetConfigCollection("golang")
}

// 获取某个appId的配置集合
func GetConfigCollection(appId string) *ConfigCollection {
	res, ok := cache[appId]
	if ok {
		return res
	}

	mux.Lock()
	defer mux.Unlock()

	//double check
	res, ok = cache[appId]
	if ok {
		return res
	}

	configApp, err := httpGetMapResp(baseUrl + "/api/getConfigApp?configAppId=" + appId)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	res = &ConfigCollection{
		appId:     appId,
		name:      configApp["name"],
		data:      map[string]*configData{},
		listeners: map[string][]ConfigChangeListener{},
	}
	res.refreshData()
	cache[appId] = res
	return res
}
