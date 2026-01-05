package initiator

import (
	"context"
	"payment-gateway/internal/service"
	"payment-gateway/internal/service/payment"
	"payment-gateway/internal/storage"
	"payment-gateway/internal/pkg/logger"
)

type Service struct {
	Payment service.PaymentService
}

func InitService(persistence Persistence, messaging storage.Messaging, log logger.Logger) Service {
	log.Info(context.Background(), "initializing payment service")

	svc := payment.NewPaymentService(persistence.Payment, messaging, log)
	log.Info(context.Background(), "payment service initialized")

	return Service{
		Payment: svc,
	}
}
