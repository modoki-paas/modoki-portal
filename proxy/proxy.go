package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (r *Request) parseURL() (*url.URL, error) {
	return url.Parse(r.URL)
}

type Proxy struct {
	clusterAddress string
	client         *http.Client
}

func NewProxy(clusterAddress string, client *http.Client) *Proxy {
	return &Proxy{
		clusterAddress: clusterAddress,
		client:         client,
	}
}

func (p *Proxy) Run(rawReq *Request) (*http.Response, error) {
	url, err := url.Parse(p.clusterAddress)

	if err != nil {
		return nil, err
	}

	innerURL, err := rawReq.parseURL()

	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, innerURL.Path)
	url.RawQuery = innerURL.Query().Encode()

	fmt.Println(url.String())

	req, err := http.NewRequest(rawReq.Method, url.String(), strings.NewReader(rawReq.Body))

	if err != nil {
		return nil, err
	}

	for k, v := range rawReq.Headers {
		req.Header.Set(k, v)
	}

	return p.client.Do(req)
}
