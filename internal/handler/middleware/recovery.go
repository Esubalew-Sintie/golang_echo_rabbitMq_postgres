package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// EchoRecovery provides panic recovery middleware
func EchoRecovery() echo.MiddlewareFunc {
	return middleware.Recover()
}
