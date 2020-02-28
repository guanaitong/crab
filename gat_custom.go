package crab

import (
	"os"
)

type GatDiscoveryClient struct {
	dnsDomainDiscoveryClient  *DnsDomainSuffixDiscoveryClient
	kubernetesDiscoveryClient *KubernetesDiscoveryClient
	inK8s                     bool
}

func NewGatDiscoveryClientInCurrentNamespace() *GatDiscoveryClient {
	return NewGatDiscoveryClient("")
}

func NewGatDiscoveryClient(namespace string) *GatDiscoveryClient {
	workEnv, workIdc := os.Getenv("WORK_ENV"), os.Getenv("WORK_IDC")
	return &GatDiscoveryClient{
		dnsDomainDiscoveryClient: &DnsDomainSuffixDiscoveryClient{
			suffix: "services." + workEnv + "." + workIdc,
		},
		kubernetesDiscoveryClient: &KubernetesDiscoveryClient{
			Namespace: namespace,
		},
		inK8s: len(os.Getenv("KUBERNETES_SERVICE_HOST")) > 0,
	}
}

func (dc *GatDiscoveryClient) GetInstances(serviceId string) (res []*ServiceInstance, e error) {
	if dc.inK8s {
		res, e = dc.kubernetesDiscoveryClient.GetInstances(serviceId)
		if e != nil && len(res) > 0 {
			return
		}
	}
	return dc.dnsDomainDiscoveryClient.GetInstances(serviceId)
}
