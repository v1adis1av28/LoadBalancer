package db

import (
	"context"
	"fmt"

	"LoadBalancer/internal/logger"

	"github.com/jackc/pgx/v5"
)

var DBConn *pgx.Conn

func InitDB(host, port, user, password, dbname string) error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	var err error
	DBConn, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		logger.Logger.Error("Unable to connect to database: %v", err)
		return err
	}

	if err := DBConn.Ping(context.Background()); err != nil {
		logger.Logger.Error("No connection with database: %v", err)
		return err
	}

	createUserTable()

	logger.Logger.Info("Successfully connected to the database")
	return nil
}

func CloseDB() {
	if DBConn != nil {
		err := DBConn.Close(context.Background())
		if err != nil {
			logger.Logger.Error("Error closing database connection: %v", err)
		} else {
			logger.Logger.Info("Database connection closed")
		}
	}
}
func createUserTable() {
	sqlStatement := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        client_id VARCHAR(255) UNIQUE NOT NULL,
        capacity INT NOT NULL,
        refill_rate INT NOT NULL
    );
    `

	_, err := DBConn.Exec(context.Background(), sqlStatement)
	if err != nil {
		logger.Logger.Error("Error creating users table: %v", err)
	}

	logger.Logger.Info("Users table created or already exists")
}
