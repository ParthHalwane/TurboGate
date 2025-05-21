package limiter

import (
	"sync"
	"time"
)

type tokenBucket struct {
	tokens         int
	lastRefillTime time.Time
	mutex          sync.Mutex
	rate           int // tokens per minute
	capacity       int
}

func newTokenBucket(rate int, capacity int) *tokenBucket {
	return &tokenBucket{
		tokens:         capacity,
		lastRefillTime: time.Now(),
		rate:           rate,
		capacity:       capacity,
	}
}

func (tb *tokenBucket) allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Minutes()
	refill := int(elapsed * float64(tb.rate))

	if refill > 0 {
		tb.tokens = min(tb.tokens+refill, tb.capacity)
		tb.lastRefillTime = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
