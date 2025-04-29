package repository

import (
	"LoadBalancer/internal/db"
	"LoadBalancer/internal/models"
	"context"
)

// AddClient добавляет нового клиента в базу данных.
func AddClient(client *models.User) error {
	query := "INSERT INTO users (client_id, capacity, refill_rate) VALUES ($1, $2, $3)"
	_, err := db.DBConn.Exec(context.Background(), query, client.ClientID, client.Capacity, client.RefillRate)
	if err != nil {
		return err
	}
	return nil
}

// DeleteClient удаляет клиента из базы данных по его ID.
func DeleteClient(clientID string) error {
	query := "DELETE FROM users WHERE client_id = $1"
	_, err := db.DBConn.Exec(context.Background(), query, clientID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllClients возвращает список всех клиентов из базы данных.
func GetAllClients() ([]models.User, error) {
	query := "SELECT client_id, capacity, refill_rate FROM users"
	rows, err := db.DBConn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ClientID, &user.Capacity, &user.RefillRate); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetClientByID возвращает клиента по его ID.
func GetClientByID(clientID string) (*models.User, error) {
	query := "SELECT client_id, capacity, refill_rate FROM users WHERE client_id = $1"
	row := db.DBConn.QueryRow(context.Background(), query, clientID)

	var user models.User
	err := row.Scan(&user.ClientID, &user.Capacity, &user.RefillRate)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
