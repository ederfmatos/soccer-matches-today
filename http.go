package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

type HttpClient struct {
	authToken string
	baseURL   string
}

func NewHttpClient(baseURL, authToken string) *HttpClient {
	return &HttpClient{baseURL: baseURL, authToken: authToken}
}

func (h HttpClient) Do(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", h.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("make request: %v", err)
	}

	req.Header.Add("x-auth-token", h.authToken)

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
