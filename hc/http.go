package hc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	u "net/url"
	"os"
	"strings"
	"time"
)

var tr = &http.Transport{
	Proxy:               http.ProxyFromEnvironment,
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	MaxIdleConns:        200,
	MaxIdleConnsPerHost: 100,
	IdleConnTimeout:     time.Duration(90) * time.Second,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
}

var client = &http.Client{
	Timeout:   time.Minute * 5, //设置一个最大超时60分钟
	Transport: tr,              // https insecure
}

func Get(url string, params map[string]interface{}) (*Resp, error) {
	var values = make(u.Values)
	if params != nil {
		for k, v := range params {
			values.Set(k, fmt.Sprint(v))
		}
	}

	if len(values) > 0 {
		if strings.Contains(url, "?") {
			url = url + "&" + values.Encode()
		} else {
			url = url + "?" + values.Encode()
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return Request(req)
}

func Delete(url string, params map[string]interface{}) (*Resp, error) {
	var values = make(u.Values)
	if params != nil {
		for k, v := range params {
			values.Set(k, fmt.Sprint(v))
		}
	}
	if len(values) > 0 {
		if strings.Contains(url, "?") {
			url = url + "&" + values.Encode()
		} else {
			url = url + "?" + values.Encode()
		}
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return Request(req)
}

func PostForm(url string, params map[string]interface{}) (*Resp, error) {
	var values = make(u.Values)
	if params != nil {
		for k, v := range params {
			values.Set(k, fmt.Sprint(v))
		}
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return Request(req)
}

func PostFormFile(url string, params map[string]interface{}) (*Resp, error) {
	// build params
	fieldName, ok := params["FormFile"].(string)
	if !ok {
		return nil, fmt.Errorf("Not found 'FormFile'")
	}
	fileName, ok := params["FileName"].(string)
	if !ok {
		return nil, fmt.Errorf("Not found 'FileName'")
	}

	// Create buffer
	buf := new(bytes.Buffer)

	// caveat IMO dont use this for large files
	// create a tmpfile and assemble your multipart from there (not tested)
	w := multipart.NewWriter(buf)

	// Create file field
	fw, err := w.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	if fd != nil {
		defer fd.Close()
	}

	// Write file field from file to upload
	_, err = io.Copy(fw, fd)
	if err != nil {
		return nil, err
	}

	// required
	if w != nil {
		w.Close()
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}

	// Important if you do not close the multipart writer you will not have a
	// terminating boundry
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Auth
	if auth, ok := params["Auth"].(map[string]string); ok {
		username, _ := auth["username"]
		passport, _ := auth["passport"]
		req.SetBasicAuth(username, passport)
	}

	// Headers
	if headers, ok := params["Headers"].(map[string]string); ok {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	return Request(req)
}

func PostJson(url string, params interface{}) (*Resp, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	return Request(req)
}

func Request(request *http.Request) (*Resp, error) {
	request.Header.Set("User-Agent", "GAT_GO_HTTP_CLIENT")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return newResp(request, response)
}
