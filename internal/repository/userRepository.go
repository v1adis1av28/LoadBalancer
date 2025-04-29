package repository

import (
	"LoadBalancer/internal/db"
	"LoadBalancer/internal/models"
)

func AddClient(client *models.User) error {
	_, err := db.DB.Exec(
		"INSERT INTO users(client_id, capacity, refill_rate) VALUES ($1, $2, $3)",
		client.ClientID, client.Capacity, client.RefillRate,
	)
	return err
}

func DeleteClient(clientID string) error {
	_, err := db.DB.Exec("DELETE FROM users WHERE client_id = $1", clientID)
	return err
}

func GetAllClients() ([]models.User, error) {
	rows, err := db.DB.Query("SELECT client_id, capacity, refill_rate FROM users")
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

func GetClientByID(clientID string) (*models.User, error) {
	row := db.DB.QueryRow("SELECT client_id, capacity, refill_rate FROM clients WHERE client_id = $1", clientID)

	var user models.User
	err := row.Scan(&user.ClientID, &user.Capacity, &user.RefillRate)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
