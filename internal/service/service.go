package service

import (
	"context"
	"payment-gateway/internal/constant/model/dto"
	"payment-gateway/internal/constant/model/response"

	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req *dto.CreatePaymentRequest) (*response.CreatePaymentResponse, error)
	GetPayment(ctx context.Context, id uuid.UUID) (*response.GetPaymentResponse, error)
	ProcessPayment(ctx context.Context, paymentID uuid.UUID) error
}
