package gconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guanaitong/crab/system"
	"io/ioutil"
	"net/http"
	u "net/url"
	"strings"
)

type gConfHttpClient struct {
	baseUrl string
}

type ConfigApp struct {
	AppId string `json:"configCollectionId"`
	Name  string `json:"name"`
}

// 获取configApp信息
func (g *gConfHttpClient) getConfigApp(appId string) *ConfigApp {
	content, err := g.getContent("/getConfigApp", map[string]string{
		"configAppId": appId,
	})
	if err != nil {
		return nil
	}
	configApp := new(ConfigApp)
	err = json.Unmarshal([]byte(content), configApp)
	if err != nil {
		return nil
	}
	return configApp
}

// 获取配置集合Key列表
func (g *gConfHttpClient) listConfigKeys(appId string) []string {
	content, err := g.getContent("/listConfigKeys", map[string]string{
		"configAppId": appId,
	})
	if err != nil {
		return []string{}
	}
	keys := make([]string, 0)
	err = json.Unmarshal([]byte(content), &keys)
	if err != nil {
		return []string{}
	}
	return keys
}

// 获取单个配置项值
func (g *gConfHttpClient) getConfig(appId, key string) string {
	content, err := g.getContent("/getConfig", map[string]string{
		"configAppId": appId,
		"key":         key,
	})
	if err != nil {
		return ""
	}
	return content
}

// 获取所有配置
func (g *gConfHttpClient) listConfigs(appId string) map[string]string {
	content, err := g.getContent("/listConfigs", map[string]string{
		"configAppId": appId,
	})
	if err != nil {
		return map[string]string{}
	}
	res := make(map[string]string)
	err = json.Unmarshal([]byte(content), &res)
	if err != nil {
		return res
	}
	return res
}

// 监听appid列表，返回需要更新的appId
func (g *gConfHttpClient) watch(configAppIds []string) []string {
	configAppIdList := strings.Join(configAppIds, ",")
	content, err := g.getContent("/watch", map[string]string{
		"configAppIdList": configAppIdList,
		"clientId":        system.GetInstanceId(),
	})
	if err != nil || content == "" {
		return []string{}
	}
	keys := make([]string, 0)
	err = json.Unmarshal([]byte(content), &keys)
	if err != nil {
		return []string{}
	}
	return keys
}

func (g *gConfHttpClient) getContent(path string, params map[string]string) (string, error) {
	url := g.baseUrl + path
	var values = make(u.Values)
	if params != nil {
		for k, v := range params {
			values.Set(k, fmt.Sprint(v))
		}
	}
	if len(values) > 0 {
		if strings.Contains(url, "?") {
			url = url + "&" + values.Encode()
		} else {
			url = url + "?" + values.Encode()
		}
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(bs), nil
	} else {
		return "", errors.New("resp status code is not 200, it it " + resp.Status + " ,url is " + url)
	}
}
