package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

type (
	HttpClient struct {
		baseURL string
		headers map[string]string
	}

	HttpClientOption func(*HttpClient)
)

func NewHttpClient(baseURL string, options ...HttpClientOption) *HttpClient {
	client := &HttpClient{
		baseURL: baseURL,
		headers: make(map[string]string),
	}
	for _, option := range options {
		option(client)
	}
	return client
}

func (h HttpClient) Do(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", h.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("make request: %v", err)
	}

	for key, value := range h.headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %v", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}

	return body, nil
}

func WithHeader(name, value string) HttpClientOption {
	return func(hc *HttpClient) {
		hc.headers[name] = value
	}
}
