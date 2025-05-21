package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify request before forwarding (optional)
	proxy.ModifyResponse = func(resp *http.Response) error {
		// You can add tracing headers or logging here
		return nil
	}

	// Optional: Error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Proxy error", http.StatusBadGateway)
	}

	return proxy
}
