package router

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
	"turbogate/config"
	"turbogate/pkg/logger"
	// other imports remain the same
)

// Remove the old tokenBucket struct, global map and mu

type TokenBucket struct {
	mu        sync.Mutex
	tokens    float64
	lastCheck time.Time
	rate      float64
	burst     float64
}

func NewTokenBucket(rate, burst float64) *TokenBucket {
	return &TokenBucket{
		tokens:    burst,
		lastCheck: time.Now(),
		rate:      rate,
		burst:     burst,
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastCheck).Seconds()
	tb.tokens += elapsed * tb.rate
	if tb.tokens > tb.burst {
		tb.tokens = tb.burst
	}
	tb.lastCheck = now

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}
	return false
}

type RateLimiter struct {
	buckets sync.Map // concurrent map[string]*TokenBucket
	rate    float64
	burst   float64
}

func NewRateLimiter(rate float64, burst float64) *RateLimiter {
	return &RateLimiter{
		rate:  rate,
		burst: burst,
	}
}

func (rl *RateLimiter) GetBucket(ip string) *TokenBucket {
	val, ok := rl.buckets.Load(ip)
	if ok {
		return val.(*TokenBucket)
	}
	tb := NewTokenBucket(rl.rate, rl.burst)
	actual, loaded := rl.buckets.LoadOrStore(ip, tb)
	if loaded {
		return actual.(*TokenBucket)
	}
	return tb
}

func (rl *RateLimiter) Allow(ip string) bool {
	bucket := rl.GetBucket(ip)
	return bucket.Allow()
}

func NewRateLimiterMiddleware(rate float64, burst float64) func(http.Handler) http.Handler {
	rl := NewRateLimiter(rate, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]

			if !rl.Allow(ip) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Keep the rest of router.go (SetupRouter, min, etc) unchanged

func SetupRouter(routes []config.Route, rateLimiter func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	for _, route := range routes {
		target, err := url.Parse(route.Target)
		if err != nil {
			logger.Error("Invalid target: " + route.Target)
			continue
		}

		transport := &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:          10000,
			MaxIdleConnsPerHost:   10000,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   20 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     false,
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = transport

		mux.Handle(route.Path, rateLimiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Strip the route prefix
			r.URL.Path = strings.TrimPrefix(r.URL.Path, route.Path)
			logger.Info("Proxying " + r.URL.Path + " → " + route.Target)
			proxy.ServeHTTP(w, r)
		})))

	}

	return mux
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
