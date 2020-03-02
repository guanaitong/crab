package k8s

type Endpoints struct {
	Kind       string                 `json:kind`
	ApiVersion string                 `json:apiVersion`
	Metadata   map[string]interface{} `json:metadata`
	Subsets    []Subset               `json:subsets`
}
type Subset struct {
	Addresses []Address `json:addresses`
	Ports     []Port    `json:ports`
}
type Address struct {
	Ip        string            `json:ip`
	NodeName  string            `json:nodeName`
	TargetRef map[string]string `json:targetRef`
}

type Port struct {
	Port     int    `json:port`
	Protocol string `json:protocol`
}
