package ms

import (
	"github.com/guanaitong/crab/gconf"
	"github.com/guanaitong/crab/ms/internal/k8s"
	"github.com/guanaitong/crab/system"
	"k8s.io/klog"
)

var (
	endpoints = map[string]string{}
)

func AddEndpointsConfig(serviceId string, host string) {
	endpoints[serviceId] = host
}

type GatDiscoveryClient struct {
	suffix                string
	domainDiscoveryClient *DnsDomainDiscoveryClient
}

func NewGatDiscoveryClient() *GatDiscoveryClient {
	return &GatDiscoveryClient{
		domainDiscoveryClient: &DnsDomainDiscoveryClient{},
		suffix:                system.GetServiceDomainSuffix(),
	}
}

func (dc *GatDiscoveryClient) GetInstances(serviceId string) (res []*ServiceInstance, e error) {
	if system.InK8s() {
		res, e = dc.getInstancesByK8s(serviceId)
		if e == nil && len(res) > 0 {
			return
		}
		if e != nil {
			klog.Errorln("cannot get instances from k8s,err:" + e.Error())
		}
	}
	host, ok := endpoints[serviceId]
	if ok {
		return dc.domainDiscoveryClient.GetInstances(host)
	}
	m := gconf.GetGlobalConfigCollection().GetConfigAsStructuredMap("service_address.properties")
	h, ok := m[serviceId]
	if ok {
		return dc.domainDiscoveryClient.GetInstances(h.AsString())
	}
	return dc.domainDiscoveryClient.GetInstances(serviceId + dc.suffix)
}

func (dc *GatDiscoveryClient) getInstancesByK8s(serviceId string) ([]*ServiceInstance, error) {
	globalConfigCollection := gconf.GetGlobalConfigCollection()
	namespaceConfig := globalConfigCollection.GetConfigAsStructuredMap("namespace.properties")
	namespace := k8s.GetCurrentNamespace()
	namespaceValue, ok := namespaceConfig[serviceId]
	if ok && namespaceValue.AsString() != "" {
		namespace = namespaceValue.AsString()
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
