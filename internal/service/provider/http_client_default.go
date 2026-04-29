//go:build !(js && wasm)

package provider

import (
	"net/http"
	"time"
)

func newHTTPClient(cfg HTTPConfig) *http.Client {
	maxIdle := cfg.MaxIdleConnsPerHost
	if maxIdle <= 0 {
		maxIdle = 4
	}

	idleTimeout := 90 * time.Second
	if cfg.IdleConnTimeout != "" {
		if d, err := time.ParseDuration(cfg.IdleConnTimeout); err == nil {
			idleTimeout = d
		}
	}

	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: maxIdle,
			IdleConnTimeout:     idleTimeout,
			DisableCompression:  false,
		},
	}
}
