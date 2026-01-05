package response

import (
	"payment-gateway/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

// CreatePaymentResponse represents the response from creating a payment
//
// swagger:model CreatePaymentResponse
type CreatePaymentResponse struct {
	// Unique payment identifier
	// example: 550e8400-e29b-41d4-a716-446655440000
	PaymentID uuid.UUID `json:"payment_id"`

	// Current payment status
	// enum: PENDING,SUCCESS,FAILED
	// example: PENDING
	Status entities.PaymentStatus `json:"status"`

	// Payment creation timestamp
	// example: 2024-01-05T10:30:00Z
	CreatedAt time.Time `json:"created_at"`
}

// GetPaymentResponse represents detailed payment information
//
// swagger:model GetPaymentResponse
type GetPaymentResponse struct {
	// Unique payment identifier
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID uuid.UUID `json:"id"`

	// Payment amount as string
	// example: 100.50
	Amount string `json:"amount"`

	// Payment currency
	// enum: ETB,USD
	// example: USD
	Currency string `json:"currency"`

	// Idempotency key used for the payment
	// example: order-12345
	IdempotencyKey string `json:"idempotency_key"`

	// Current payment status
	// enum: PENDING,SUCCESS,FAILED
	// example: SUCCESS
	Status entities.PaymentStatus `json:"status"`

	// Payment creation timestamp
	// example: 2024-01-05T10:30:00Z
	CreatedAt time.Time `json:"created_at"`

	// Payment processing completion timestamp (null if still pending)
	// example: 2024-01-05T10:30:05Z
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
}
