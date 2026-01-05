package core

import (
	"context"
	"payment-gateway/internal/pkg/logger"
)

type Core struct {
	Logger logger.Logger
}

func NewCore(logger logger.Logger) *Core {
	return &Core{
		Logger: logger,
	}
}

func (c *Core) Health(ctx context.Context) error {
	return nil
}
