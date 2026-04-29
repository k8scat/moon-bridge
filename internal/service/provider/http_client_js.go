//go:build js && wasm

package provider

import (
	"net/http"
)

func newHTTPClient(cfg HTTPConfig) *http.Client {
	// js/wasm (e.g. Cloudflare Workers) uses the browser fetch API;
	// http.Transport fields are not available.
	return http.DefaultClient
}
