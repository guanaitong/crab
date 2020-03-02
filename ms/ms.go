package ms

import "net/http"

// Version # of crab
const Version = "1.0.0-SNAPSHOT"

// ServiceInstance just represents a application unit,it represents a pod in k8s
type ServiceInstance struct {
	ServiceId  string
	InstanceId string
	Host       string
	Ip         string
	Port       int
	Status     status
}

type status struct {
	NetFailed bool
}

type ClientBuilder struct {
	serviceId       string
	discoveryClient DiscoveryClient
	loadBalance     LoadBalance
	httpClient      *http.Client
	Debug           bool
}

// New method creates a new crab service client builder.
func New(serviceId string) *ClientBuilder {
	if serviceId == "" {
		panic("serviceId cannot be empty")
	}
	return &ClientBuilder{
		serviceId: serviceId,
	}
}

func (cb *ClientBuilder) SetDiscoveryClient(client DiscoveryClient) *ClientBuilder {
	cb.discoveryClient = client
	return cb
}

func (cb *ClientBuilder) SetHttpClient(httpClient *http.Client) *ClientBuilder {
	cb.httpClient = httpClient
	return cb
}

func (cb *ClientBuilder) SetLoadBalance(lb LoadBalance) *ClientBuilder {
	cb.loadBalance = lb
	return cb
}

func (cb *ClientBuilder) SetDebug(debug bool) *ClientBuilder {
	cb.Debug = debug
	return cb
}

func (cb *ClientBuilder) Build() *ServiceClient {
	loadBalance := cb.loadBalance
	if loadBalance == nil {
		loadBalance = &RandomLoadBalance{}
	}
	discoveryClient := cb.discoveryClient
	if discoveryClient == nil {
		discoveryClient = new(DnsDomainDiscoveryClient)
	}
	httpClient := cb.httpClient
	if httpClient == nil {
		httpClient = globalHttpClient
	}
	serviceClient := &ServiceClient{
		serviceId:   cb.serviceId,
		cache:       globalCache.newServiceCache(cb.serviceId, discoveryClient),
		loadBalance: loadBalance,
		httpClient:  httpClient,
		Debug:       cb.Debug,
	}
	return serviceClient
}
