package initiator

import (
	"net/http"
	"payment-gateway/internal/constant/errors"
	"payment-gateway/internal/constant/model/response"
	"payment-gateway/internal/glue/routing/payment"
	"payment-gateway/internal/handler/middleware"
	"payment-gateway/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, handler Handler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	e.GET("/health", handler.Payment.HealthCheck)

	e.POST("/api/v1/auth/login", func(c echo.Context) error {
		userID := "user123"
		username := "testuser"
		email := "test@example.com"

		token, err := authMiddleware.GenerateToken(userID, username, email)
		if err != nil {
			return c.JSON(errors.GetHTTPStatus(errors.ErrInternalServerError), response.Response{
				Status: errors.GetHTTPStatus(errors.ErrInternalServerError),
				Error: &response.DetailedErrorResponse{
					Type:    "server_error",
					Message: errors.GetUserFriendlyMessage(errors.ErrInternalServerError),
				},
			})
		}

		authData := map[string]interface{}{
			"token":    token,
			"user_id":  userID,
			"username": username,
			"email":    email,
		}

		return c.JSON(http.StatusOK, response.Response{
			Status:  http.StatusOK,
			Message: "Authentication successful",
			Data:    authData,
		})
	})

	paymentGroup := e.Group("/api/v1/payments")
	payment.InitPaymentRoutes(paymentGroup, handler.Payment, authMiddleware)
}
