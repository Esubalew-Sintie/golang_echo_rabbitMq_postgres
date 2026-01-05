package initiator

import (
	"context"
	"payment-gateway/internal/pkg/logger"
)

func InitLogger(ctx context.Context) logger.Logger {
	logger := logger.InitLogger()
	logger.Info(ctx, "logger initialized")
	return logger
}
