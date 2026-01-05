package payment

import (
	"payment-gateway/internal/handler"
	"payment-gateway/internal/handler/middleware"

	"github.com/labstack/echo/v4"
)

func InitPaymentRoutes(paymentGroup *echo.Group, paymentHandler handler.PaymentHandler, authMiddleware *middleware.AuthMiddleware) {
	// paymentGroup.POST("", paymentHandler.CreatePayment, echo.WrapMiddleware(authMiddleware.AuthenticateToken))
	// paymentGroup.GET("/:id", paymentHandler.GetPayment, echo.WrapMiddleware(authMiddleware.AuthenticateToken))
	paymentGroup.POST("", paymentHandler.CreatePayment)
	paymentGroup.GET("/:id", paymentHandler.GetPayment)

}
