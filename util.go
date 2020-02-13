package crab

import (
	"net/http"
	"reflect"
	"strings"
)

// isStringEmpty method tells whether given string is empty or not
func isStringEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// detectContentType method is used to figure out `Request.body` content type for request header
func detectContentType(body interface{}) string {
	contentType := plainTextType
	kind := kindOf(body)
	switch kind {
	case reflect.Struct, reflect.Map:
		contentType = jsonContentType
	case reflect.String:
		contentType = plainTextType
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			contentType = jsonContentType
		}
	}

	return contentType
}

// isJSONType method is to check JSON content type or not
func isJSONType(ct string) bool {
	return jsonCheck.MatchString(ct)
}

func isPayloadSupported(m string) bool {
	return m == http.MethodPost || m == http.MethodPut || m == http.MethodPatch
}

func getPointer(v interface{}) interface{} {
	vv := valueOf(v)
	if vv.Kind() == reflect.Ptr {
		return v
	}
	return reflect.New(vv.Type()).Interface()
}

func valueOf(i interface{}) reflect.Value {
	return reflect.ValueOf(i)
}

func kindOf(v interface{}) reflect.Kind {
	return typeOf(v).Kind()
}

func typeOf(i interface{}) reflect.Type {
	return indirect(valueOf(i)).Type()
}

func indirect(v reflect.Value) reflect.Value {
	return reflect.Indirect(v)
}
