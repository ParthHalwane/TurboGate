package router

import (
	"net/http"
	"strings"

	"TurboGate/config"
	"TurboGate/internal/limiter"
	"TurboGate/internal/proxy"
)

func NewRouter(cfg *config.Config, lim *limiter.IPRateLimiter) http.Handler {
	mux := http.NewServeMux()

	for _, route := range cfg.Routes {
		target := route.Target
		path := route.Path

		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			if !lim.Allow(clientIP) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Strip prefix for cleaner forwarding
			r.URL.Path = strings.TrimPrefix(r.URL.Path, path)
			proxy.NewReverseProxy(target).ServeHTTP(w, r)
		})
	}

	return mux
}
