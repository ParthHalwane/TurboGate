package limiter

import (
	"net"
	"sync"
)

type IPRateLimiter struct {
	clients map[string]*tokenBucket
	mutex   sync.Mutex
	rate    int
	cap     int
}

func NewIPRateLimiter(rate int, capacity int) *IPRateLimiter {
	return &IPRateLimiter{
		clients: make(map[string]*tokenBucket),
		rate:    rate,
		cap:     capacity,
	}
}

func (i *IPRateLimiter) getBucket(ip string) *tokenBucket {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if bucket, exists := i.clients[ip]; exists {
		return bucket
	}

	bucket := newTokenBucket(i.rate, i.cap)
	i.clients[ip] = bucket
	return bucket
}

func (i *IPRateLimiter) Allow(remoteAddr string) bool {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		ip = remoteAddr
	}
	bucket := i.getBucket(ip)
	return bucket.allow()
}
