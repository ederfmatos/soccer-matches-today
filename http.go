package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	HttpClient struct {
		baseURL string
		headers map[string]string
		client  *http.Client
	}

	HttpClientOption func(*HttpClient)
)

func NewHttpClient(baseURL string, options ...HttpClientOption) *HttpClient {
	client := &HttpClient{
		baseURL: baseURL,
		headers: make(map[string]string),
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
	for _, option := range options {
		option(client)
	}
	return client
}

func (h HttpClient) Get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", h.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("make request: %v", err)
	}

	for key, value := range h.headers {
		req.Header.Add(key, value)
	}

	res, err := h.client.Do(req)
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

func (h HttpClient) Post(path string, body interface{}) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal data: %v", err)
	}

	req, err := http.NewRequest("POST", h.baseURL+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("make request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	for key, value := range h.headers {
		req.Header.Add(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status: %d - [%s]", resp.StatusCode, responseBody)
	}

	return responseBody, nil
}

func WithHeader(name, value string) HttpClientOption {
	return func(hc *HttpClient) {
		hc.headers[name] = value
	}
}
