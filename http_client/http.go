package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	MaxIdleConns        = 100
	MaxIdleConnsPerHost = 100
	IdleConnTimeout     = 90 * time.Second
	KeepAlive           = 30 * time.Second
	Timeout             = 30 * time.Second
)

type Options struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	Timeout             time.Duration
	KeepAlive           time.Duration
}

type Client http.Client

func New(usePool bool, op Options) *Client {
	var client = new(Client)
	if usePool {
		if op.MaxIdleConns == 0 {
			op.MaxIdleConns = MaxIdleConns
		}
		if op.MaxIdleConnsPerHost == 0 {
			op.MaxIdleConnsPerHost = MaxIdleConnsPerHost
		}
		if op.IdleConnTimeout == 0 {
			op.IdleConnTimeout = IdleConnTimeout
		}
		if op.KeepAlive == 0 {
			op.IdleConnTimeout = KeepAlive
		}
		if op.Timeout == 0 {
			op.Timeout = Timeout
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   op.Timeout,
				KeepAlive: op.KeepAlive,
			}).DialContext,
			MaxIdleConns:        op.MaxIdleConns,
			MaxIdleConnsPerHost: op.MaxIdleConnsPerHost,
			IdleConnTimeout:     op.IdleConnTimeout,
		}
	}
	return client
}

func (c *Client) ReqJson(method string, url string, header http.Header,
	InputData, outputData interface{}) (int, error) {

	body, _ := json.Marshal(InputData)
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header = header
	response, err := (*http.Client)(c).Do(req)
	if err != nil {
		return 0, err
	}
	if response.StatusCode != http.StatusOK {
		return response.StatusCode, fmt.Errorf("[%d] http request failed", response.StatusCode)
	}

	by, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return http.StatusOK, err
	}

	err = json.Unmarshal(by, outputData)
	response.Body.Close()

	return http.StatusOK, err
}
