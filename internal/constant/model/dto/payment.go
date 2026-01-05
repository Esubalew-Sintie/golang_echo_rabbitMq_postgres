package dto

import (
	"payment-gateway/internal/domain/entities"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CreatePaymentRequest represents a request to create a new payment
//
// swagger:model CreatePaymentRequest
type CreatePaymentRequest struct {
	// Payment amount (must be greater than 0)
	// required: true
	// example: 100.50
	Amount float64 `json:"amount"`

	// Payment currency (ETB or USD)
	// required: true
	// enum: ETB,USD
	// example: USD
	Currency string `json:"currency"`

	// Unique idempotency key to prevent duplicate payments
	// required: true
	// example: order-12345
	IdempotencyKey string `json:"idempotency_key"`
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
