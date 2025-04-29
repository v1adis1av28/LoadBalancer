package main

import (
	"LoadBalancer/internal/balancer"
	"LoadBalancer/internal/config"
	"LoadBalancer/internal/db"
	"LoadBalancer/internal/handlers"
	"LoadBalancer/internal/logger"
	"LoadBalancer/internal/proxy"
	"LoadBalancer/internal/rateLimiter"
	"encoding/json"
	"net/http"
)

type Message struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func main() {
	logger.SetupLogger("app.log")

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		logger.Logger.Error("Error while loading config", err.Error())
	}

	b := balancer.NewLoadBalancer(cfg.Backends)
	p := proxy.NewProxy(b)
	rl := rateLimiter.NewRateLimiter()
	//	ch := &handlers.ClientsHandler{Limiter: rl}
	db.InitDB("postgres", "5432:5432", "postgres", "postgres", "balancer")
	mux := http.NewServeMux()

	// API для клиентов
	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.AddClient(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/clients/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handlers.DeleteClient(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.RemoteAddr
		if !rl.Allow(clientID) {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(Message{Code: "429", Message: "Rate limit exceeded"})
			return
		}
		p.Serve(w, r)
	})

	logger.Logger.Info("server start on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		logger.Logger.Error("Error starting server: %v", err)
	}
}
