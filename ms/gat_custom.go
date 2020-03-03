package ms

import (
	"github.com/guanaitong/crab/gconf"
	"github.com/guanaitong/crab/ms/internal/k8s"
	"github.com/guanaitong/crab/system"
)

type GatDiscoveryClient struct {
	dnsDomainDiscoveryClient *DnsDomainSuffixDiscoveryClient
}

func NewGatDiscoveryClient() *GatDiscoveryClient {
	return &GatDiscoveryClient{
		dnsDomainDiscoveryClient: &DnsDomainSuffixDiscoveryClient{
			Suffix: system.GetServiceDomainSuffix(),
		},
	}
}

func (dc *GatDiscoveryClient) GetInstances(serviceId string) (res []*ServiceInstance, e error) {
	if system.InK8s() {
		res, e = dc.getInstancesByK8s(serviceId)
		if e != nil && len(res) > 0 {
			return
		}
	}
	return dc.dnsDomainDiscoveryClient.GetInstances(serviceId)
}

func (dc *GatDiscoveryClient) getInstancesByK8s(serviceId string) ([]*ServiceInstance, error) {
	namespaceConfig := make(map[string]string)
	globalConfigCollection := gconf.GetGlobalConfigCollection()
	globalConfigCollection.GetConfigAsBean("namespace.properties", namespaceConfig)
	namespace, ok := namespaceConfig[serviceId]
	if !ok {
		namespace = k8s.GetCurrentNamespace()
	}

	token := globalConfigCollection.GetConfig("crab_k8s_token")
	endpoints, e := k8s.GetEndpoints(serviceId, namespace, token)
	if e != nil {
		return nil, e
	}
	var res []*ServiceInstance
	for _, subset := range endpoints.Subsets {
		var port = 80
		if len(subset.Ports) == 1 {
			port = subset.Ports[0].Port
		} else {
			// 默认取TCP协议的，这里需要扩展
			for _, portObject := range subset.Ports {
				if portObject.Protocol == "TCP" {
					port = portObject.Port
				}
			}
		}
		for _, address := range subset.Addresses {
			instance := ServiceInstance{
				ServiceId:  serviceId,
				InstanceId: address.TargetRef["uid"],
				Host:       address.Ip,
				Ip:         address.Ip,
				Port:       port,
			}
			res = append(res, &instance)
		}

	}
	return res, nil
}
