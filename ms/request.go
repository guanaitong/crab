package ms

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	client           *ServiceClient
	serviceInstance  *ServiceInstance
	secure           bool
	method           string
	path             string
	pathParams       map[string]string
	queryParam       url.Values
	url              string
	header           http.Header
	formData         url.Values
	body             interface{}
	bodyBytes        []byte
	UserInfo         *User
	Token            string
	Time             time.Time
	Result           interface{}
	Error            interface{}
	RawRequest       *http.Request
	ctx              context.Context
	notParseResponse bool
}

// User type is to hold an username and password information
type User struct {
	Username, Password string
}

func (r *Request) SetSecure(secure bool) *Request {
	r.secure = secure
	return r
}

// SetHeader method is to set a single header field and its value in the current request.
func (r *Request) SetHeader(header, value string) *Request {
	r.header.Set(header, value)
	return r
}

// SetHeaders method sets multiple headers field and its values at one go in the current request.
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for h, v := range headers {
		r.SetHeader(h, v)
	}
	return r
}

// SetQueryParam method sets single parameter and its value in the current request.
func (r *Request) SetQueryParam(param, value string) *Request {
	r.queryParam.Set(param, value)
	return r
}

// SetQueryParams method sets multiple parameters and its values at one go in the current request.
// It will be formed as query string for the request.
func (r *Request) SetQueryParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetQueryParam(p, v)
	}
	return r
}

// SetQueryString method provides ability to use string as an input to set URL query string for the request.
func (r *Request) SetQueryString(query string) *Request {
	params, err := url.ParseQuery(strings.TrimSpace(query))
	if err == nil {
		for p, v := range params {
			for _, pv := range v {
				r.queryParam.Add(p, pv)
			}
		}
	} else {
		//TODO
	}
	return r
}

// SetFormData method sets Form parameters and their values in the current request.
func (r *Request) SetFormData(data map[string]string) *Request {
	for k, v := range data {
		r.formData.Set(k, v)
	}
	return r
}

func (r *Request) SetDoNotParseResponse(parse bool) *Request {
	r.notParseResponse = parse
	return r
}

func (r *Request) SetBody(body interface{}) *Request {
	r.body = body
	return r
}

func (r *Request) SetResult(res interface{}) *Request {
	r.Result = getPointer(res)
	return r
}

func (r *Request) SetError(err interface{}) *Request {
	r.Error = getPointer(err)
	return r
}

func (r *Request) SetPathParams(params map[string]string) *Request {
	for p, v := range params {
		r.pathParams[p] = v
	}
	return r
}

func (r *Request) SetPath(path string) *Request {
	r.path = path
	return r
}

func (r *Request) SetBasicAuth(username, password string) *Request {
	r.UserInfo = &User{Username: username, Password: password}
	return r
}

func (r *Request) SetAuthToken(token string) *Request {
	r.Token = token
	return r
}

func (r *Request) Get() (*Response, error) {
	r.method = http.MethodGet
	return r.execute()
}

func (r *Request) Post() (*Response, error) {
	r.method = http.MethodPost
	return r.execute()
}

func (r *Request) Head() (*Response, error) {
	r.method = http.MethodHead
	return r.execute()
}

func (r *Request) Put() (*Response, error) {
	r.method = http.MethodPut
	return r.execute()
}

func (r *Request) Delete() (*Response, error) {
	r.method = http.MethodDelete
	return r.execute()
}

func (r *Request) Options() (*Response, error) {
	r.method = http.MethodOptions
	return r.execute()
}

func (r *Request) Patch() (*Response, error) {
	r.method = http.MethodPatch
	return r.execute()
}

func (r *Request) execute() (*Response, error) {
	return r.client.doReq(r)
}
