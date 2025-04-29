package config

import (
	"LoadBalancer/internal/logger"
	"encoding/json"
	"os"
)

type Config struct {
	Port     string   `json:"port"`
	Backends []string `json:"backends"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		logger.Logger.Error("Error while opening file:", "file", path)
		return nil, err
	}
	defer file.Close()

	logger.Logger.Info("Config file was upload", path)

	var config Config

	decode := json.NewDecoder(file)
	if err := decode.Decode(&config); err != nil {
		logger.Logger.Error("Error while decoding config.json", "path", path)
		return nil, err
	}

	return &config, nil
}
