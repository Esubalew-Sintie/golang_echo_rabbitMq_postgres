package initiator

import (
	"context"
	"payment-gateway/internal/pkg/config"
	"payment-gateway/internal/pkg/logger"
)

func InitConfig(ctx context.Context, logger logger.Logger) (*config.Config, error) {
	logger.Info(ctx, "initializing configuration")
	config, err := config.InitConfig()
	if err != nil {
		logger.Error(ctx, "failed to initialize configuration: %v", err)
		return nil, err
	}
	logger.Info(ctx, "configuration initialized successfully")
	return config, nil
}
