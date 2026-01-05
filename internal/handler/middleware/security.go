package middleware

import (
	"github.com/labstack/echo/v4"
)

// EchoSecurityHeaders adds security headers to responses
func EchoSecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Prevent clickjacking
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// Prevent MIME type sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// Enable XSS protection
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			// Referrer policy
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Content Security Policy - basic policy for API
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")

			return next(c)
		}
	}
}
