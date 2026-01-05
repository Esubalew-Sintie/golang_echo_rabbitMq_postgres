package payment

import (
	"net/http"
	"payment-gateway/internal/constant/errors"
	"payment-gateway/internal/constant/model/dto"
	"payment-gateway/internal/constant/model/response"
	"payment-gateway/internal/handler"
	"payment-gateway/internal/pkg/logger"
	"payment-gateway/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	service service.PaymentService
	log     logger.Logger
}

func NewPaymentHandler(service service.PaymentService, log logger.Logger) handler.PaymentHandler {
	return &PaymentHandler{
		service: service,
		log:     log,
	}
}

func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.CreatePaymentRequest

	if err := c.Bind(&req); err != nil {
		h.log.Error(ctx, "Failed to decode request: %v", err)
		return response.SendValidationError(c, errors.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		h.log.Error(ctx, "Request validation failed: %v", err)
		fieldErrors := response.ErrorFields(err)
		return response.SendValidationErrorResponse(c, "Request validation failed. Please check your input.", fieldErrors)
	}

	result, err := h.service.CreatePayment(ctx, &req)
	if err != nil {
		h.log.Error(ctx, "Failed to create payment: %v", err)
		return response.SendErrorResponse(c, err)
	}

	return response.SendSuccessResponse(c, http.StatusCreated, "Payment created successfully", result, nil)
}

func (h *PaymentHandler) GetPayment(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	paymentID, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Error(ctx, "Invalid payment ID format: %s, error: %v", idStr, err)
		return response.SendValidationError(c, errors.ErrInvalidID)
	}

	payment, err := h.service.GetPayment(ctx, paymentID)
	if err != nil {
		h.log.Error(ctx, "Failed to get payment: %s, error: %v", paymentID.String(), err)
		return response.SendErrorResponse(c, err)
	}

	return response.SendSuccessResponse(c, http.StatusOK, "Payment retrieved successfully", payment, nil)
}

func (h *PaymentHandler) HealthCheck(c echo.Context) error {
	healthData := map[string]interface{}{
		"status":  "healthy",
		"service": "payment-gateway",
		"version": "1.0.0",
	}
	return response.SendSuccessResponse(c, http.StatusOK, "Service is healthy", healthData, nil)
}
