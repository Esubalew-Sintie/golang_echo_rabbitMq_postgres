package dto

import (
	"payment-gateway/internal/domain/entities"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreatePaymentRequest struct {
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotency_key"`
}

func (r CreatePaymentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Amount,
			validation.Required.Error("amount is required"),
			validation.Min(0.01).Error("amount must be greater than zero"),
		),
		validation.Field(&r.Currency,
			validation.Required.Error("currency is required"),
			validation.In("ETB", "USD").Error("currency must be either ETB or USD"),
		),
		validation.Field(&r.IdempotencyKey,
			validation.Required.Error("idempotency key is required"),
			validation.Length(1, 255).Error("idempotency key must be between 1 and 255 characters"),
		),
	)
}

// ProcessPaymentMessage represents the message sent to RabbitMQ for processing
type ProcessPaymentMessage struct {
	PaymentID string `json:"payment_id"`
}

// ToPaymentEntity converts the DTO to a Payment entity
func (r *CreatePaymentRequest) ToPaymentEntity() (*entities.Payment, error) {
	// This will be implemented when we create the service layer
	return nil, nil
}
