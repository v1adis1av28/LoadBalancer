package db

import (
	"LoadBalancer/internal/logger"
	"database/sql"
	"fmt"
)

var DB *sql.DB

func InitDB(host, port, user, password, dbname string) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Logger.Error("Error while connecting db", "host", host, "port", port, "DB name", dbname)
		return err
	}

	if err = DB.Ping(); err != nil {
		logger.Logger.Error("No connection with db", "host", host, "port", port, "DB name", dbname)
		return err
	}

	query := `
    CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(255) UNIQUE NOT NULL,
    capacity INT NOT NULL,
    refill_rate INT NOT NULL
  );
  `
	_, err = DB.Exec(query)
	return err
}
