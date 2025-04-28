package balancer

import (
	"LoadBalancer/internal/logger"
	"sync"
)

type LoadBalancer struct {
	Index    int
	Backends []string
	mu       sync.Mutex
}

// Init func for Balancer
func NewLoadBalancer(backends []string) *LoadBalancer {
	return &LoadBalancer{Index: 0, Backends: backends}
}

func NextBackend(lb *LoadBalancer) string {

	lb.mu.Lock()
	defer lb.mu.Unlock()

	//RoundRobin implementation
	backend := lb.Backends[lb.Index]
	lb.Index = (lb.Index + 1) % len(lb.Backends)
	logger.Logger.Info("Current backend:", "backend", backend)
	return backend
}
