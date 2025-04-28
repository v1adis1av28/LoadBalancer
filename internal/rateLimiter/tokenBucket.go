package rateLimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	RefilRate  int
	mu         sync.Mutex
	lastRefill time.Time
	Capacity   int
	Tokens     int
}

func (rl *RateLimiter) refillTokens() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		rl.mu.Lock()
		for _, bucket := range rl.buckets {
			bucket.mu.Lock()
			now := time.Now()
			sub := now.Sub(bucket.lastRefill).Seconds()
			TokensAddCount := int(sub * float64(bucket.RefilRate))
			if TokensAddCount > 0 {
				bucket.Tokens += TokensAddCount
				if bucket.Tokens > bucket.Capacity {
					bucket.Tokens = bucket.Capacity
				}
				bucket.lastRefill = now
			}
			bucket.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}
