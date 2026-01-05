package payment

import (
	"context"
	customErrors "errors"
	"math/rand"
	"payment-gateway/internal/constant/errors"
	"payment-gateway/internal/constant/model/dto"
	"payment-gateway/internal/constant/model/response"
	"payment-gateway/internal/domain/entities"
	"payment-gateway/internal/pkg/logger"
	"payment-gateway/internal/storage"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentService struct {
	persistence storage.Persistence
	messaging   storage.Messaging
	log         logger.Logger
}

func NewPaymentService(
	persistence storage.Persistence,
	messaging storage.Messaging,
	log logger.Logger,
) *PaymentService {
	return &PaymentService{
		persistence: persistence,
		messaging:   messaging,
		log:         log,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *dto.CreatePaymentRequest) (*response.CreatePaymentResponse, error) {
	s.log.Info(ctx, "Creating payment with idempotency guarantee: reference=%s, currency=%s", req.IdempotencyKey, req.Currency)

	tx, err := s.persistence.BeginTx(ctx)
	if err != nil {
		s.log.Error(ctx, "Failed to begin database transaction: %v", err)
		return nil, errors.ErrDatabaseOperationFailed
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				s.log.Error(ctx, "Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	payment := &entities.Payment{
		ID:             uuid.New(),
		Amount:         decimal.NewFromFloat(req.Amount),
		Currency:       req.Currency,
		IdempotencyKey: req.IdempotencyKey,
		Status:         entities.PaymentStatusPending,
		CreatedAt:      time.Now(),
	}
	existingPayment, getErr := s.persistence.GetPaymentByReference(ctx, req.IdempotencyKey)
	if getErr != nil {
		s.log.Error(ctx, "Failed to fetch existing payment after unique constraint error: %v", getErr)
		if !customErrors.Is(getErr, errors.ErrPaymentNotFound) {
			return nil, errors.ErrDuplicateTransaction
		}
	}
	if existingPayment != nil {
		s.log.Info(ctx, "Returning existing payment for idempotency key %s: %s", req.IdempotencyKey, existingPayment.ID.String())
		return nil, errors.ErrDuplicateTransaction
	}

	err = s.persistence.CreatePaymentWithTx(ctx, tx, payment)
	if err != nil {
		s.log.Error(ctx, "Failed to create payment in transaction: %v", err)
		return nil, err
	}

	err = s.messaging.PublishPaymentProcessing(ctx, payment.ID)
	if err != nil {
		s.log.Error(ctx, "Failed to publish payment processing message: %s, error: %v", payment.ID.String(), err)
		return nil, errors.ErrMessagePublishingFailed
	}

	if err = tx.Commit(); err != nil {
		s.log.Error(ctx, "Failed to commit transaction: %v", err)
		return nil, errors.ErrDatabaseOperationFailed
	}

	s.log.Info(ctx, "Payment created successfully with idempotency guarantee: %s, reference: %s", payment.ID.String(), req.IdempotencyKey)

	return &response.CreatePaymentResponse{
		PaymentID: payment.ID,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
	}, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*response.GetPaymentResponse, error) {
	s.log.Info(ctx, "Retrieving payment: %s", id.String())

	payment, err := s.persistence.GetPaymentByID(ctx, id)
	if err != nil {
		s.log.Error(ctx, "Failed to retrieve payment: %s, error: %v", id.String(), err)
		return nil, err
	}

	return &response.GetPaymentResponse{
		ID:             payment.ID,
		Amount:         payment.Amount.String(),
		Currency:       payment.Currency,
		IdempotencyKey: payment.IdempotencyKey,
		Status:         entities.PaymentStatus(payment.Status),
		CreatedAt:      payment.CreatedAt,
		ProcessedAt:    payment.ProcessedAt,
	}, nil
}

func (s *PaymentService) ProcessPayment(ctx context.Context, paymentID uuid.UUID) error {
	s.log.Info(ctx, "Processing payment: %s", paymentID.String())

	newStatus := s.simulatePaymentProcessing()
	s.log.Info(ctx, "Payment processing simulated: %s, target_status: %s", paymentID.String(), string(newStatus))

	err := s.persistence.ProcessPaymentWithTransaction(ctx, paymentID, newStatus)
	if err != nil {
		s.log.Error(ctx, "Failed to process payment transactionally: %s, error: %v", paymentID.String(), err)

		switch err {
		case errors.ErrPaymentNotFound:
			return errors.ErrPaymentNotFound
		case errors.ErrPaymentAlreadyProcessed:
			s.log.Info(ctx, "Payment already processed: %s", paymentID.String())
			return errors.ErrPaymentAlreadyProcessed
		case errors.ErrInvalidPaymentStatus:
			return errors.ErrInvalidPaymentStatus
		default:
			return errors.ErrPaymentInitiationFailed
		}
	}

	s.log.Info(ctx, "Payment processed successfully: %s, final_status: %s", paymentID.String(), string(newStatus))
	return nil
}

func (s *PaymentService) simulatePaymentProcessing() entities.PaymentStatus {
	if rand.Float32() < 0.5 {
		return entities.PaymentStatusSuccess
	}
	return entities.PaymentStatusFailed
}
