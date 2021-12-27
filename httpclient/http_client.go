package httpclient

import (
	"net/http"
	"time"
)

func NewDefault() *http.Client {
	tr := &http.Transport{
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConnsPerHost: 1000,
		TLSHandshakeTimeout: 1 * time.Second,
	}

	return &http.Client{Transport: tr, Timeout: 3 * time.Second}
}
