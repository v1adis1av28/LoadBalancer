package rateLimiter

import (
	"LoadBalancer/internal/logger"
	"sync"
	"time"
)

type RateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*TokenBucket),
	}
	//Start refiling in gorutine
	go rl.refillTokens()
	logger.Logger.Info("Create rate limiter and start refill user tokens")
	return rl
}

func (rl *RateLimiter) Allow(clientId string) bool {
	rl.mu.Lock()
	tokenBucket, exists := rl.buckets[clientId]
	if !exists {
		tokenBucket = &TokenBucket{
			RefilRate:  1,
			Capacity:   10,
			Tokens:     10,
			lastRefill: time.Now(),
		}
		rl.buckets[clientId] = tokenBucket
	}
	rl.mu.Unlock()

	tokenBucket.mu.Lock()
	defer tokenBucket.mu.Unlock()

	if tokenBucket.Tokens > 0 {
		tokenBucket.Tokens--
		logger.Logger.Info("Allow user access", "user", clientId)
		return true
	}
	logger.Logger.Info("Doesn`t allow access for user", "user", clientId)
	return false
}
