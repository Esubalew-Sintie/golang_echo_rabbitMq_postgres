package storage

import (
	"context"
	"payment-gateway/internal/constant/model/dto"
	"payment-gateway/internal/domain/entities"

	"github.com/google/uuid"
)

type Persistence interface {
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status entities.PaymentStatus) error

	BeginTx(ctx context.Context) (Transaction, error)
	ProcessPaymentWithTransaction(ctx context.Context, paymentID uuid.UUID, status entities.PaymentStatus) error
	CreatePayment(ctx context.Context, payment *entities.Payment) error
	GetPaymentByIdempotencyKey(ctx context.Context, idempotencyKey string) (*entities.Payment, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}

type Messaging interface {
	PublishPaymentProcessing(ctx context.Context, paymentID uuid.UUID) error
	ConsumePaymentProcessing() (<-chan dto.ProcessPaymentMessage, error)
	Close() error
}
