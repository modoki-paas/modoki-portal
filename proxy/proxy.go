package main

import (
	"bytes"
	"net/http"
	"net/url"
)

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
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
	url.RawPath = innerURL.RawPath

	req, err := http.NewRequest(rawReq.Method, url.String(), bytes.NewReader(rawReq.Body))

	if err != nil {
		return nil, err
	}

	for k, v := range rawReq.Headers {
		req.Header.Set(k, v)
	}

	return p.client.Do(req)
}
