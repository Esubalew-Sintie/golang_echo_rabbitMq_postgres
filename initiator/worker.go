package initiator

import (
	"context"
	"payment-gateway/internal/pkg/logger"
	"payment-gateway/internal/rabbitmq"
	"payment-gateway/internal/storage"
)

func InitWorker(service Service, messaging storage.Messaging, logger logger.Logger) *rabbitmq.PaymentWorker {
	logger.Info(context.Background(), "initializing payment worker")

	worker := rabbitmq.NewPaymentWorker(service.Payment, messaging, logger)
	logger.Info(context.Background(), "payment worker initialized")

	return worker
}
