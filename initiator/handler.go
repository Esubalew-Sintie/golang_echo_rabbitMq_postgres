package initiator

import (
	"payment-gateway/internal/handler"
	httppayment "payment-gateway/internal/handler/rest/http/payment"
	"payment-gateway/internal/pkg/logger"
)

type Handler struct {
	Payment handler.PaymentHandler
}

func InitHandler(service Service, log logger.Logger) Handler {
	hdlr := httppayment.NewPaymentHandler(service.Payment, log)
	return Handler{
		Payment: hdlr,
	}
}
