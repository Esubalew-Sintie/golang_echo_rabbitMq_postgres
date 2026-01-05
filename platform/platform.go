package platform

import (
	"payment-gateway/internal/pkg/logger"
	"payment-gateway/platform/core"
)

// Platform holds all platform services
type Platform struct {
	Core   *core.Core
	Logger logger.Logger
}

// InitPlatform initializes the platform services
func InitPlatform(log logger.Logger) *Platform {

	core := core.NewCore(log)

	return &Platform{
		Core:   core,
		Logger: log,
	}
}
