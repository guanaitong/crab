package k8s

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/guanaitong/crab/json"
	"io/ioutil"
	"net"
	"os"
	"time"
)
import "net/http"

var k8sServiceUrl string
var currentNamespace string

var tr = &http.Transport{
	Proxy:               http.ProxyFromEnvironment,
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	MaxIdleConns:        200,
	MaxIdleConnsPerHost: 100,
	IdleConnTimeout:     time.Duration(90) * time.Second,
	DialContext: (&net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 60 * time.Second,
	}).DialContext,
}

var k8sHttpClient = &http.Client{
	Timeout:   time.Second * 300, //设置一个最大超时300秒
	Transport: tr,                // https insecure
}

// 这里不适用k8s的client-go,它太庞大,我们只需要一个“获取service endPoint信息”的api，所以自己撸
func GetEndpoints(serviceId, namespace, token string) (*Endpoints, error) {
	url := GetK8sServiceUrl() + "/api/v1/namespaces/" + namespace + "/endpoints/" + serviceId
	endpoints := &Endpoints{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	e := doReq(req, endpoints)
	if e != nil {
		return nil, fmt.Errorf("get k8s endpoints error:%w", e)
	}
	return endpoints, nil
}

func GetK8sServiceUrl() string {
	if k8sServiceUrl == "" {
		host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
		if len(host) == 0 || len(port) == 0 {
			fmt.Println("application does not run in kubernetes")
			panic("application does not run in kubernetes")
		}
		k8sServiceUrl = "https://" + host + ":" + port
	}
	return k8sServiceUrl
}

// https://stackoverflow.com/questions/46046110/how-to-get-the-current-namespace-in-a-pod
func GetCurrentNamespace() string {
	if currentNamespace == "" {
		data, e := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if e != nil {
			fmt.Println("application cannot read file /var/run/secrets/kubernetes.io/serviceaccount/namespace")
			panic("application cannot read file /var/run/secrets/kubernetes.io/serviceaccount/namespace")
		}
		currentNamespace = string(data)
	}
	return currentNamespace
}

func doReq(req *http.Request, v interface{}) error {
	resp, err := k8sHttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bs, v)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("resp status code is not 200, it it " + resp.Status)
	}
}
