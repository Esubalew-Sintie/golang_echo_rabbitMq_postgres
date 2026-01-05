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
	// Swagger routes (before security middleware)
	e.GET("/swagger.json", func(c echo.Context) error {
		return c.File("docs/swagger.json")
	})

	e.GET("/swagger", func(c echo.Context) error {
		html := `<!DOCTYPE html>
<html>
<head>
  <title>Payment Gateway API - Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui.css" />
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: '/swagger.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`
		return c.HTML(200, html)
	})

	e.GET("/swagger/", func(c echo.Context) error {
		return c.Redirect(301, "/swagger")
	})

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
