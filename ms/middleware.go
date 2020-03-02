package ms

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type buildRequestFunc func(*Request) error

type buildResponseFunc func(*Response) error

var crabBuildRequestFunc = []buildRequestFunc{
	parseRequestURL,
	parseRequestHeader,
	parseRequestBody,
	createHTTPRequest,
	addCredentials,
}

func buildHttpRequest(r *Request) error {
	for _, f := range crabBuildRequestFunc {
		if err := f(r); err != nil {
			return err
		}
	}
	return nil
}

var crabBuildResponseFunc = []buildResponseFunc{
	parseResponseBody,
}

func buildHttpResponse(r *Response) error {
	for _, f := range crabBuildResponseFunc {
		if err := f(r); err != nil {
			return err
		}
	}
	return nil
}

func parseRequestURL(r *Request) error {
	var urlBuilder strings.Builder
	if r.secure {
		urlBuilder.WriteString("https://")
	} else {
		urlBuilder.WriteString("http://")
	}
	urlBuilder.WriteString(r.serviceInstance.Ip)
	if !(r.serviceInstance.Port == 0) {
		urlBuilder.WriteString(":")
		urlBuilder.WriteString(strconv.Itoa(r.serviceInstance.Port))
	}
	if !strings.HasPrefix(r.path, "/") {
		urlBuilder.WriteString("/")
	}
	var path = r.path
	if len(r.pathParams) > 0 {
		for p, v := range r.pathParams {
			path = strings.Replace(path, "{"+p+"}", url.PathEscape(v), -1)
		}
	}

	urlBuilder.WriteString(path)

	if len(r.queryParam) > 0 {
		var values = make(url.Values)
		for k, v := range r.queryParam {
			for _, iv := range v {
				values.Add(k, iv)
			}
		}
		urlBuilder.WriteString("?")
		urlBuilder.WriteString(values.Encode())
	}
	reqUrl := urlBuilder.String()

	// Parsing request URL
	_, err := url.Parse(reqUrl)
	if err != nil {
		return err
	}
	r.url = reqUrl
	return nil
}

func parseRequestHeader(r *Request) error {
	hdr := make(http.Header)
	for k := range r.header {
		hdr[k] = append(hdr[k], r.header[k]...)
	}

	if isStringEmpty(hdr.Get(hdrUserAgentKey)) {
		hdr.Set(hdrUserAgentKey, hdrUserAgentValue)
	}

	ct := hdr.Get(hdrContentTypeKey)
	if isStringEmpty(hdr.Get(hdrAcceptKey)) && !isStringEmpty(ct) && (isJSONType(ct)) {
		hdr.Set(hdrAcceptKey, hdr.Get(hdrContentTypeKey))
	}

	r.header = hdr
	return nil
}

func parseRequestBody(r *Request) (err error) {
	if isPayloadSupported(r.method) {
		// Handling Form Data
		if len(r.formData) > 0 {
			handleFormData(r)
		}
		// Handling Request body
		if r.body != nil {
			handleContentType(r)
			if err = handleRequestBody(r); err != nil {
				return
			}
		}
	}
	return
}

func handleFormData(r *Request) {
	formData := url.Values{}

	for k, v := range r.formData {
		for _, iv := range v {
			formData.Add(k, iv)
		}
	}

	r.bodyBytes = []byte(formData.Encode())
	r.header.Set(hdrContentTypeKey, formContentType)
}

func handleContentType(r *Request) {
	contentType := r.header.Get(hdrContentTypeKey)
	if isStringEmpty(contentType) {
		contentType = detectContentType(r.body)
		r.header.Set(hdrContentTypeKey, contentType)
	}
}

func handleRequestBody(r *Request) (err error) {
	var bodyBytes []byte = nil
	contentType := r.header.Get(hdrContentTypeKey)
	kind := kindOf(r.body)

	if reader, ok := r.body.(io.Reader); ok {
		bodyBytes, err = ioutil.ReadAll(reader)
		if err != nil {
			return
		}
	} else if b, ok := r.body.([]byte); ok {
		bodyBytes = b
	} else if s, ok := r.body.(string); ok {
		bodyBytes = []byte(s)
	} else if isJSONType(contentType) &&
		(kind == reflect.Struct || kind == reflect.Map || kind == reflect.Slice) {
		bodyBytes, err = json.Marshal(r.body)
		if err != nil {
			return
		}
	}

	if bodyBytes == nil {
		err = errors.New("unsupported 'body' type/value")
		return
	}
	r.bodyBytes = bodyBytes
	return
}

func createHTTPRequest(r *Request) (err error) {
	if len(r.bodyBytes) == 0 {
		r.RawRequest, err = http.NewRequest(r.method, r.url, nil)
	} else {
		r.RawRequest, err = http.NewRequest(r.method, r.url, bytes.NewBuffer(r.bodyBytes))
	}

	if err != nil {
		return
	}

	// Add headers into http request
	r.RawRequest.Header = r.header
	r.RawRequest.Host = r.serviceInstance.Host

	// Use context if it was specified
	if r.ctx != nil {
		r.RawRequest = r.RawRequest.WithContext(r.ctx)
	}

	// assign get body func for the underlying raw request instance
	r.RawRequest.GetBody = func() (io.ReadCloser, error) {
		// If r.bodyBuf present, return the copy
		if len(r.bodyBytes) > 0 {
			return ioutil.NopCloser(bytes.NewReader(r.bodyBytes)), nil
		}
		return nil, nil
	}
	return nil
}

func addCredentials(r *Request) error {
	// Basic Auth
	if r.UserInfo != nil { // takes precedence
		r.RawRequest.SetBasicAuth(r.UserInfo.Username, r.UserInfo.Password)
	}
	// Token Auth
	if !isStringEmpty(r.Token) { // takes precedence
		r.RawRequest.Header.Set(hdrAuthorizationKey, "Bearer "+r.Token)
	}
	return nil
}

func parseResponseBody(res *Response) (err error) {
	if res.StatusCode() == http.StatusNoContent {
		return
	}
	// Handles only JSON type
	ct := res.Header().Get(hdrContentTypeKey)
	if isJSONType(ct) {
		// HTTP status code > 199 and < 300, considered as Result
		if res.IsSuccess() {
			res.Request.Error = nil
			if res.Request.Result != nil {
				err = json.Unmarshal(res.body, res.Request.Result)
				return
			}
		}

		// HTTP status code > 399, considered as Error
		if res.IsError() {
			if res.Request.Error != nil {
				err = json.Unmarshal(res.body, res.Request.Error)
			}
		}
	}

	return
}
