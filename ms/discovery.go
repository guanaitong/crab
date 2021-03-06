package ms

import (
	"fmt"
	"github.com/guanaitong/crab/ms/internal/k8s"
	"net"
)

// Represents read operations commonly available to discovery services
type DiscoveryClient interface {
	GetInstances(serviceId string) ([]*ServiceInstance, error)
}

// base on k8s
type KubernetesDiscoveryClient struct {
	Namespace string
	Token     string
}

func (dc *KubernetesDiscoveryClient) GetInstances(serviceId string) ([]*ServiceInstance, error) {
	namespace := dc.Namespace
	if namespace == "" {
		namespace = k8s.GetCurrentNamespace()
	}
	endpoints, e := k8s.GetEndpoints(serviceId, namespace, dc.Token)
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

type DnsDomainDiscoveryClient struct {
}

func (dc *DnsDomainDiscoveryClient) GetInstances(domainAsServiceId string) ([]*ServiceInstance, error) {
	addrs, err := net.LookupHost(domainAsServiceId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var res []*ServiceInstance
	for _, addr := range addrs {
		instance := ServiceInstance{
			ServiceId:  domainAsServiceId,
			InstanceId: addr,
			Host:       domainAsServiceId,
			Ip:         addr,
		}
		res = append(res, &instance)
	}
	return res, nil
}
