package handlers

import (
	"LoadBalancer/internal/logger"
	"LoadBalancer/internal/rateLimiter"
	"encoding/json"
	"net/http"
	"strings"
)

// struct for post requests
type UserRequest struct {
	ClientID   string `json:"client_id"`
	Capacity   int    `json:"capacity"`
	RefillRate int    `json:"rate_per_sec"`
}

type UserHandler struct {
	Limiter *rateLimiter.RateLimiter
}

func (h *UserHandler) AddClient(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error on request body!")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.Limiter.AddUser(req.ClientID, req.Capacity, req.RefillRate)
	logger.Logger.Info("User was added", req.ClientID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Client added"))
}

func (h *UserHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimPrefix(r.URL.Path, "/clients/")
	if clientID == "" {
		logger.Logger.Error("Error while deleting user! Empty clientID")
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	h.Limiter.DeleteUser(clientID)
	logger.Logger.Info("User was deleted", "UserID", clientID)
	w.Write([]byte("Client deleted"))
}

//TODO add sql implementation
