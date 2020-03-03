package ms

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	hdrUserAgentKey       = http.CanonicalHeaderKey("User-Agent")
	hdrAcceptKey          = http.CanonicalHeaderKey("Accept")
	hdrContentTypeKey     = http.CanonicalHeaderKey("Content-Type")
	hdrContentLengthKey   = http.CanonicalHeaderKey("Content-Length")
	hdrContentEncodingKey = http.CanonicalHeaderKey("Content-Encoding")
	hdrAuthorizationKey   = http.CanonicalHeaderKey("Authorization")

	plainTextType   = "text/plain; charset=utf-8"
	jsonContentType = "application/json"
	formContentType = "application/x-www-form-urlencoded"

	jsonCheck = regexp.MustCompile(`(?i:(application|text)/(json|.*\+json|json\-.*)(;|$))`)
	xmlCheck  = regexp.MustCompile(`(?i:(application|text)/(xml|.*\+xml)(;|$))`)

	hdrUserAgentValue = "crab-client/" + Version + " (https://github.com/guanaitong/crab)"
	bufPool           = &sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}
)

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

var globalHttpClient = &http.Client{
	Timeout:   time.Second * 300, //设置一个最大超时300秒
	Transport: tr,                // https insecure
}

type ServiceClient struct {
	serviceId       string
	DiscoveryClient DiscoveryClient
	cache           *serviceCache
	loadBalance     LoadBalance
	httpClient      *http.Client
	Debug           bool
}

func (client *ServiceClient) R() *Request {
	return &Request{
		client:     client,
		queryParam: url.Values{},
		formData:   url.Values{},
		header:     http.Header{},
		pathParams: map[string]string{},
	}
}

func (client *ServiceClient) doReq(request *Request) (*Response, error) {
	for i := 0; ; i++ {
		resp, retry, err := client.doOneReq(request)
		if retry && i < 3 {
			continue
		}
		return resp, err
	}
}

func (client *ServiceClient) doOneReq(request *Request) (*Response, bool, error) {
	serviceInstances := client.cache.GetUpdatedInstances()
	serviceInstance := client.loadBalance.DoSelect(serviceInstances)
	if serviceInstance == nil {
		return nil, false, errors.New("no instance found")
	}
	request.serviceInstance = serviceInstance
	err := buildHttpRequest(request)
	if err != nil {
		return nil, false, fmt.Errorf("illegal request %w", err)
	}
	request.Time = time.Now()
	httpResponse, err := client.httpClient.Do(request.RawRequest)
	if err != nil {
		errString := err.Error()
		if strings.Contains(errString, "dial tcp") {
			serviceInstance.Status.NetFailed = true
		}
		return nil, true, fmt.Errorf("dial tcp error %w", err)
	}
	endTime := time.Now()

	response := &Response{
		Request:     request,
		RawResponse: httpResponse,
		receivedAt:  endTime,
	}
	if err != nil {
		//value, ok := err.(url.Error)
		return response, false, fmt.Errorf("request error %w", err)
	}
	if request.notParseResponse {
		return response, false, nil
	}

	defer httpResponse.Body.Close()

	if response.body, err = ioutil.ReadAll(httpResponse.Body); err != nil {
		return response, false, fmt.Errorf("read body error %w", err)
	}
	if err = buildHttpResponse(response); err != nil {
		return response, false, fmt.Errorf("parse response error %w", err)
	}
	return response, false, nil
}
