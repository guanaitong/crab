package gconf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func httpGetMapResp(url string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var m = make(map[string]string)
		err = json.Unmarshal(bs, &m)
		if err != nil {
			return nil, err
		}
		return m, nil
	} else {
		return nil, errors.New("resp status code is not 200, it it " + resp.Status + " ,baseUrl is " + url)
	}
}

func httpGetListResp(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var m []string
		err = json.Unmarshal(bs, &m)
		if err != nil {
			return nil, err
		}
		return m, nil
	} else {
		return nil, errors.New("resp status code is not 200, it it " + resp.Status + " ,baseUrl is " + url)
	}
}
