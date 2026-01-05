package initiator

import (
	"context"
	"os"
	"payment-gateway/internal/pkg/logger"
	"payment-gateway/internal/storage"
	"payment-gateway/internal/storage/messaging"
)

func InitMessageBroker(logger logger.Logger) (storage.Messaging, error) {
	logger.Info(context.Background(), "initializing RabbitMQ message broker")
	host := getEnvOrDefault("RABBITMQ_HOST", "localhost")
	port := getEnvOrDefault("RABBITMQ_PORT", "5672")
	user := getEnvOrDefault("RABBITMQ_USER", "guest")
	password := getEnvOrDefault("RABBITMQ_PASSWORD", "guest")
	vhost := getEnvOrDefault("RABBITMQ_VHOST", "/")

	logger.Info(context.Background(), "connecting to RabbitMQ at %s:%s", host, port)

	broker, err := messaging.NewRabbitMQBroker(host, port, user, password, vhost, logger)
	if err != nil {
		logger.Error(context.Background(), "failed to initialize RabbitMQ broker: %v", err)
		return nil, err
	}

	logger.Info(context.Background(), "RabbitMQ message broker initialized")
	return broker, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
