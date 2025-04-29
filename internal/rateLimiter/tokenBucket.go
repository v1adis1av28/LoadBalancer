package rateLimiter

import (

	"LoadBalancer/internal/logger"
	"encoding/json"
	"os"

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

func getDefaulBucket() *TokenBucket {
	file, err := os.Open("bucket.json")
	if err != nil {
		logger.Logger.Error("Error while openning file bucket.json")
		return nil
	}
	defer file.Close()
	logger.Logger.Info("Default bucket file was upload")
	var bucket TokenBucket

	decode := json.NewDecoder(file)
	if err := decode.Decode(&bucket); err != nil {
		logger.Logger.Error("Error while decoding bucket.json")
		return nil
	}
	bucket.lastRefill = time.Now()
	return &bucket
}
