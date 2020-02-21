package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type BaseResp struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func ProxyRequest(req *http.Request, host string, queries map[string]string) (*http.Request, error) {
	url := req.URL
	url.Host = host
	url.Scheme = "http"

	q := url.Query()
	for k, v := range queries {
		q.Add(k, v)
	}
	url.RawQuery = q.Encode()

	proxyReq, err := http.NewRequest(req.Method, url.String(), req.Body)
	if err != nil {
		return nil, err
	}
	proxyReq.Header.Set("Host", req.Host)
	proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}
	return proxyReq, nil
}

func GetBodyAsByte(resp *http.Response) ([]byte, error) {
	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	} else {
		return nil, nil
	}
	return ioutil.ReadAll(resp.Body)
}

func GetBodyAsString(resp *http.Response) (string, error) {
	bs, err := GetBodyAsByte(resp)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func GetBodyAsStruct(resp *http.Response, s interface{}) error {
	bs, err := GetBodyAsByte(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, s)
}

