package main

import (
	"LoadBalancer/internal/balancer"
	"LoadBalancer/internal/config"
	"LoadBalancer/internal/logger"
	"LoadBalancer/internal/proxy"
	"LoadBalancer/internal/rateLimiter"
	"net/http"
)

func StartApp() {
	//Setuping logger
	logger.SetupLogger("app.log")
	// Configuration init
	cfg, err := config.LoadConfig("G:/LoadBalancer/pkg/models/config.json")
	if err != nil {
		logger.Logger.Error("Error while loading config", err.Error())
	}

	b := balancer.NewLoadBalancer(cfg.Backends)
	p := proxy.NewProxy(b)
	rl := rateLimiter.NewRateLimiter()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.RemoteAddr
		if !rl.Allow(clientID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		p.Serve(w, r)
	})

	// Starting server
	logger.Logger.Info("server start on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		logger.Logger.Error("Error starting server: %v", err)
	}

}
