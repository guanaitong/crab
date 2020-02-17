# crab
关爱通go语言微服务架构，面向k8s

## 架构

整体架构比较类似于springcloud+kubernetes：https://github.com/spring-cloud/spring-cloud-kubernetes

### **服务端架构**

- 服务端协议为http+json的方式，使用gin框架
- 应用使用deployment部署在k8s之中，配置k8s service

服务端比较简单，开发和部署都很容易，类似于一个http服务。

### **客户端架构**

整个微服务框架其实采用了客户端负载均衡的方式，所以整个框架主要集中在客户端。

#### 服务发现

服务发现定义为`DiscoveryClient`

框架实现了三种服务发现，应用运行在k8s中和非k8s中(本地跑代码)

- `KubernetesDiscoveryClient`：运行在k8s，那么可以此实现。它会通过k8s的api(endpoints接口)，获取生产列表
- `DnsDomainDiscoveryClient`：基于dns的负载均衡。它会根据dns解析解析得到IP列表
- `DnsDomainSuffixDiscoveryClient`：本质也是基于dns，只是可以使用统一的dns后缀。一般企业生产环境内网，都会有一个统一内部域名后缀。

用户可以自定义`DiscoveryClient`实现

#### 服务负载均衡

负载均衡为`LoadBalance`，默认为`RandomLoadBalance`。

用户可以自定义`LoadBalance`实现。

#### 服务实例缓存

并不是每次调用远程服务的时候，都会调用`DiscoveryClient`实时取一下实例列表。这样子做性能有很大的问题，而且在高qps的情况下，容易把服务注册中心(比如k8s api server)弄挂掉。而且，服务注册中心的实例状态，更新也是有延迟的，并不说明获取的列表全部是可用的。

我们实现了服务实例的缓存，有以下几种情况会更新缓存：

1. 初始化的时候
2. 当服务缓存里，可用服务个数为空时，会刷新
3. 定时的全局任务，会刷新所有服务的服务实例

#### 服务健康机制

服务注册中心里提供的服务实例列表，只能保证其在一定时间范围内是有效的。而客户端对服务实例的调用是否正常，才是服务实例健康状况的最直接有效证明。

目前，我们简单实现了，在网络不通的情况下，将实例标记为不可用状态，后续的调度中，就会忽略不可用实例。

#### 调用重试

目前，只有在网络不通的情况下，才会简单的重试三次。

## 核心api

对用户的api，参考了resty的实现：https://github.com/go-resty/resty

### serviceClient

通过`clientBuilder`创建serviceClient，可以设置定义的`DiscoveryClient`、`LoadBalance`、`http.Client`等。`serviceClient`只有一个`R()`方法，使用它会得到一个`Request`对象

### Request

request对象用于定义http request请求。

**特别声明一点**：request对象无法定义请求的`baseURL`和header里的`HOST`字段。因为它们都是交给框架处理了，中间涉及到服务调用的逻辑。



举例：

```
m := map[string]interface{}{}
	resp, err := crab.
		New("httpbin.org").
		Build().
		R().
		SetPath("/post").
		SetQueryParam("k1", "v1").
		SetFormData(map[string]string{
			"userId":       "sample@sample.com",
			"subAccountId": "100002",
		}).
		SetResult(&m).
		Post()
fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse body: %v", resp)
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
```

可以采用builder模式，一路下来

更多例子可以看[crab_test.go](crab_test.go)

## 使用规范

- 服务提供方，除了需要写服务端代码，还需要编写客户端的sdk。
- 一个服务对应一个serviceClient对象

### 项目结构

1. 核心业务代码都放在pkg包里
2. apis也位于pkg之中，里面有api.go和xx.go。api.go为使用crab定义的客户端sdk，xx.go中为struct对象列表。