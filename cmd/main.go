package main

import (
	"LoadBalancer/internal/balancer"
	"LoadBalancer/internal/config"
	"LoadBalancer/internal/logger"
	"LoadBalancer/internal/proxy"
	"LoadBalancer/internal/rateLimiter"
	"encoding/json"
	"net/http"
)

type Message struct {
	Code    string `json : "code"`
	Message string `json : "message"`
}

func main() {
	//Starting application
	//Setuping logger
	logger.SetupLogger("app.log")
	// Configuration init
	cfg, err := config.LoadConfig("config.json")
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
			obj, err := json.Marshal(Message{Code: "429", Message: "Rate limit exceed!"})
			if err != nil {
				logger.Logger.Error("Error while decoding message")
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(obj)
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
