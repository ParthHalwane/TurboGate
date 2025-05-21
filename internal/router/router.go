package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"TurboGate/config"
	"TurboGate/internal/limiter"
	"TurboGate/internal/proxy"
	"TurboGate/pkg/logger"
)

func NewRouter(cfg *config.Config, lim *limiter.IPRateLimiter) http.Handler {
	mux := http.NewServeMux()

	for _, route := range cfg.Routes {
		target := route.Upstream
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

type tokenBucket struct {
	tokens         int
	maxTokens      int
	refillInterval time.Duration
	lastRefill     time.Time
	mu             sync.Mutex
}

var rateLimiters = make(map[string]*tokenBucket)
var mu sync.Mutex

// NewRateLimiterMiddleware creates a middleware with rate limiting per IP
func NewRateLimiterMiddleware(maxTokens int, refillRatePerSec int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]

			mu.Lock()
			limiter, exists := rateLimiters[ip]
			if !exists {
				limiter = &tokenBucket{
					tokens:         maxTokens,
					maxTokens:      maxTokens,
					refillInterval: time.Second / time.Duration(refillRatePerSec),
					lastRefill:     time.Now(),
				}
				rateLimiters[ip] = limiter
			}
			mu.Unlock()

			limiter.mu.Lock()
			defer limiter.mu.Unlock()

			now := time.Now()
			elapsed := now.Sub(limiter.lastRefill)
			tokensToAdd := int(elapsed / limiter.refillInterval)
			if tokensToAdd > 0 {
				limiter.tokens = min(limiter.maxTokens, limiter.tokens+tokensToAdd)
				limiter.lastRefill = now
			}

			if limiter.tokens <= 0 {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			limiter.tokens--
			next.ServeHTTP(w, r)
		})
	}
}

// SetupRouter takes parsed routes and rate limiter, returns a handler
func SetupRouter(routes []config.Route, rateLimiter func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	for _, route := range routes {
		target, err := url.Parse(route.Upstream)
		if err != nil {
			logger.Error("Invalid upstream: " + route.Upstream)
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		handler := rateLimiter(proxy)

		mux.Handle(route.Path, handler)
	}

	return mux
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
