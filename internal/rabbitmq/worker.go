package rabbitmq

import (
	"context"
	"payment-gateway/internal/constant/errors"
	"payment-gateway/internal/constant/model/dto"
	"payment-gateway/internal/service"
	"payment-gateway/internal/storage"
	"payment-gateway/internal/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type PaymentWorker struct {
	paymentService service.PaymentService
	messaging      storage.Messaging
	logger         logger.Logger
}

func NewPaymentWorker(paymentService service.PaymentService, messaging storage.Messaging, logger logger.Logger) *PaymentWorker {
	return &PaymentWorker{
		paymentService: paymentService,
		messaging:      messaging,
		logger:         logger,
	}
}

func (w *PaymentWorker) Start(ctx context.Context) error {
	w.logger.Info(ctx, "Starting payment processing worker")

	messages, err := w.messaging.ConsumePaymentProcessing()
	if err != nil {
		w.logger.Error(ctx, "Failed to get message channel: %v", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			w.logger.Info(ctx, "Payment worker stopping due to context cancellation")
			return nil
		case msg, ok := <-messages:
			if !ok {
				w.logger.Info(ctx, "Message channel closed, stopping worker")
				return nil
			}

			w.processMessage(ctx, msg)
		}
	}
}

func (w *PaymentWorker) processMessage(ctx context.Context, msg dto.ProcessPaymentMessage) {
	w.logger.Info(ctx, "Processing payment message: %s", msg.PaymentID)

	paymentID, err := uuid.Parse(msg.PaymentID)
	if err != nil {
		w.logger.Error(ctx, "Invalid payment ID in message: %s, error: %v", msg.PaymentID, err)
		return
	}

	processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err = w.paymentService.ProcessPayment(processCtx, paymentID)
	if err != nil {
		w.logger.Error(ctx, "Failed to process payment: %s, error: %v", paymentID.String(), err)
		return
	}

	w.logger.Info(ctx, "Payment processed successfully: %s", paymentID.String())
}

func (w *PaymentWorker) isRecoverableError(err error) bool {
	switch err {
	case errors.ErrPaymentAlreadyProcessed:
		return false
	case errors.ErrPaymentNotFound:
		return false
	default:
		return true
	}
}
