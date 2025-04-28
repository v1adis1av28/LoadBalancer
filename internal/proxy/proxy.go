package proxy

import (
	"LoadBalancer/internal/balancer"
	"LoadBalancer/internal/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	balancer *balancer.LoadBalancer
}

func NewProxy(lb *balancer.LoadBalancer) *Proxy {
	return &Proxy{balancer: lb}
}

func (p *Proxy) Serve(w http.ResponseWriter, rq *http.Request) {
	backendUrl, _ := url.Parse(p.balancer.NextBackend())
	proxy := httputil.NewSingleHostReverseProxy(backendUrl)

	proxy.ErrorHandler = func(w http.ResponseWriter, rq *http.Request, e error) {
		http.Error(w, "Backend server error", http.StatusBadGateway)
		logger.Logger.Error("Backend server error calling in:", "method", "Serve")
	}
	logger.Logger.Info("Serving request", "backendUrl", backendUrl)
	proxy.ServeHTTP(w, rq)
}
