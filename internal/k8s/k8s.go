package k8s

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)
import "net/http"

var k8sServiceUrl string
var currentNamespace string

// 这里不适用k8s的client-go,它太庞大,我们只需要一个“获取service endPoint信息”的api，所以自己撸
func GetEndpoints(serviceId string) (*Endpoints, error) {
	url := getK8sServiceUrl() + "/api/v1/namespaces/" + getCurrentNamespace() + "/endpoints/" + serviceId
	endpoints := &Endpoints{}
	e := httpGet(url, endpoints)
	if e != nil {
		return nil, fmt.Errorf("get k8s endpoints error:%w", e)
	}
	return endpoints, nil
}

func getK8sServiceUrl() string {
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
func getCurrentNamespace() string {
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

func httpGet(url string, v interface{}) error {
	resp, err := http.Get(url)
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
		return errors.New("resp status code is not 200, it it " + resp.Status + " ,url is " + url)
	}
}
