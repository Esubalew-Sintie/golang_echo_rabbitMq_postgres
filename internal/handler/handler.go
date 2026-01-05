package handler

import "github.com/labstack/echo/v4"

type PaymentHandler interface {
	CreatePayment(c echo.Context) error
	GetPayment(c echo.Context) error
	HealthCheck(c echo.Context) error
}
