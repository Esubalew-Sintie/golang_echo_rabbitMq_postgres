package response

import (
	"payment-gateway/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type CreatePaymentResponse struct {
	PaymentID uuid.UUID              `json:"payment_id"`
	Status    entities.PaymentStatus `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
}

type GetPaymentResponse struct {
	ID             uuid.UUID              `json:"id"`
	Amount         string                 `json:"amount"`
	Currency       string                 `json:"currency"`
	IdempotencyKey string                 `json:"idempotency_key"`
	Status         entities.PaymentStatus `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
	ProcessedAt    *time.Time             `json:"processed_at,omitempty"`
}
