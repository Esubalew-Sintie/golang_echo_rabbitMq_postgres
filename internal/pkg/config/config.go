package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string
	ServerAddr       string
	JwtSecretKey     string
}

func InitConfig() (*Config, error) {
	portStr := os.Getenv("DATABASE_PORT")
	port := 5432 // default
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
	}

	config := &Config{
		DatabaseHost:     getEnvOrDefault("DATABASE_HOST", "localhost"),
		DatabasePort:     port,
		DatabaseUser:     getEnvOrDefault("DATABASE_USER", "postgres"),
		DatabasePassword: getEnvOrDefault("DATABASE_PASSWORD", "12345678"),
		DatabaseName:     getEnvOrDefault("DATABASE_NAME", "ussd"),
		DatabaseSSLMode:  getEnvOrDefault("DATABASE_SSL_MODE", "disable"),
		ServerAddr:       getEnvOrDefault("SERVER_ADDR", ":8080"),
		JwtSecretKey:     getEnvOrDefault("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
	}
	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
