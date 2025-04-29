package handlers

import (
	"LoadBalancer/internal/logger"
	"LoadBalancer/internal/models"
	"LoadBalancer/internal/rateLimiter"
	"LoadBalancer/internal/repository"
	"encoding/json"
	"net/http"
	"strings"
)

type ClientRequest struct {
	ClientID   string `json:"client_id"`
	Capacity   int    `json:"capacity"`
	RefillRate int    `json:"rate_per_sec"`
}

type ClientsHandler struct {
	Limiter *rateLimiter.RateLimiter
}

func AddClient(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Logger.Error("Invalid json body on adding client")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := repository.AddClient(&user); err != nil {
		http.Error(w, "Failed to add client", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Client created"))
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimPrefix(r.URL.Path, "/clients/")
	if clientID == "" {
		logger.Logger.Error("Error on client id", "client_id", clientID)
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteClient(clientID); err != nil {
		logger.Logger.Error("Error while deleting client", "client_id", clientID)
		http.Error(w, "Failed to delete client", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Client deleted"))
}
