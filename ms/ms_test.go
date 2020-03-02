package ms_test

import (
	"fmt"
	"github.com/guanaitong/crab/ms"
	"testing"
)

func TestGet(t *testing.T) {
	resp, err := ms.New("httpbin.org").
		Build().
		R().
		SetPath("/get").
		SetQueryParam("k1", "v1").
		SetQueryParam("k2", "v2").
		SetHeader("gat", "great").
		Get()
	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse body: %v", resp)
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
}

func TestPostBody(t *testing.T) {
	resp, err := ms.New("httpbin.org").
		Build().
		R().
		SetPath("/post").
		SetBody(&ms.User{Username: "xx", Password: "yy"}).
		Post()
	print(t, resp, err)

}

func TestPostFormBody(t *testing.T) {
	m := map[string]interface{}{}
	resp, err := ms.New("httpbin.org").
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
	print(t, resp, err)
}

func print(t *testing.T, resp *ms.Response, err error) {
	if err != nil {
		t.Fail()
	}
	if resp.IsError() {
		t.Fail()
	}
	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse body: %v", resp)
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
}
