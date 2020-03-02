package hc

import (
	"encoding/json"
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
)

type Resp struct {
	raw      *http.Response
	bodyData []byte
}

func newResp(req *http.Request, resp *http.Response) (*Resp, error) {
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		klog.Warning("resp status code is " + resp.Status + " ,url is " + req.URL.String() + " ,body is " + string(d))
	}
	res := Resp{
		raw:      resp,
		bodyData: d,
	}
	return &res, nil
}

func (resp *Resp) StatusCode() int {
	return resp.raw.StatusCode
}

func (resp *Resp) Status() string {
	return resp.raw.Status
}
func (resp *Resp) Header() http.Header {
	return resp.raw.Header
}
func (resp *Resp) AsString() string {
	return string(resp.bodyData)
}

func (resp *Resp) AsJson(v interface{}) error {
	return json.Unmarshal(resp.bodyData, v)
}
