package repository

import (
	"LoadBalancer/internal/db"
	"LoadBalancer/internal/models"
	"context"
)

func AddClient(client *models.User) error {
	query := "INSERT INTO users (client_id, capacity, refill_rate) VALUES ($1, $2, $3)"
	_, err := db.DBConn.Exec(context.Background(), query, client.ClientID, client.Capacity, client.RefillRate)
	if err != nil {
		return err
	}
	return nil
}

func DeleteClient(clientID string) error {
	query := "DELETE FROM users WHERE client_id = $1"
	_, err := db.DBConn.Exec(context.Background(), query, clientID)
	if err != nil {
		return err
	}
	return nil
}
