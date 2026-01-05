package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
	PaymentStatusFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	Amount         decimal.Decimal `json:"amount" db:"amount"`
	Currency       string          `json:"currency" db:"currency"`
	IdempotencyKey string          `json:"idempotency_key" db:"idempotency_key"`
	Status         PaymentStatus   `json:"status" db:"status"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	ProcessedAt    *time.Time      `json:"processed_at,omitempty" db:"processed_at"`
}

func (p *Payment) IsTerminalState() bool {
	return p.Status == PaymentStatusSuccess || p.Status == PaymentStatusFailed
}

func (p *Payment) CanBeProcessed() bool {
	return p.Status == PaymentStatusPending
}
