package http

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"
)

// NewClient new client with timeouts
func NewClient() *http.Client {
	var netClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
	return netClient
}

// Get http GET with client that has timeouts
func Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return NewClient().Do(req)
}

// Post http POST with client that has timeouts
func Post(url string, contentType string, body string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if err != nil {
		return nil, err
	}
	return NewClient().Do(req)
}

// Do request with client that has timeouts
func Do(method, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return NewClient().Do(req)
}
