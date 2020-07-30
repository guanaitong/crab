package gconf

import (
	"github.com/guanaitong/crab/system"
	"sync"
	"time"
)

var ds *dataStore

func init() {
	domain := ""
	if system.InK8s() {
		domain = "gconf.kube-system"
	} else {
		domain = "gconf" + system.GetServiceDomainSuffix()
	}
	client := &gConfHttpClient{
		baseUrl: "http://" + domain + "/api",
	}
	ds = &dataStore{
		dataCache: map[string]*ConfigCollection{},
		client:    client,
		mux:       new(sync.Mutex),
	}
	ds.startBackgroundTask()
}

type dataStore struct {
	dataCache map[string]*ConfigCollection
	client    *gConfHttpClient
	mux       *sync.Mutex
}

func (ds *dataStore) startBackgroundTask() {
	go func() {
		for {
			if len(ds.dataCache) == 0 {
				time.Sleep(time.Second * 2)
				continue
			}
			var appIdList []string
			for k := range ds.dataCache {
				appIdList = append(appIdList, k)
			}
			needChangeAppIdList := ds.client.watch(appIdList)

			for _, appId := range needChangeAppIdList {
				ds.dataCache[appId].refreshData(ds.client)
			}
		}
	}()
}

func (ds *dataStore) getConfigCollection(appId string) *ConfigCollection {
	res, ok := ds.dataCache[appId]
	if ok {
		return res
	}

	ds.mux.Lock()
	defer ds.mux.Unlock()

	//double check
	res, ok = ds.dataCache[appId]
	if ok {
		return res
	}

	configApp := ds.client.getConfigApp(appId)
	if configApp == nil {
		return nil
	}

	res = &ConfigCollection{
		appId:     appId,
		name:      configApp.Name,
		data:      map[string]*Value{},
		listeners: map[string][]ConfigChangeListener{},
	}
	res.refreshData(ds.client)
	ds.dataCache[appId] = res
	return res
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
	return ds.getConfigCollection(appId)
}
