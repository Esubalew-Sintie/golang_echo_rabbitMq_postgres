package initiator

import (
	"context"
	"payment-gateway/internal/pkg/logger"
	"payment-gateway/internal/storage"
	"payment-gateway/internal/storage/persistence"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Persistence struct {
	Payment storage.Persistence
}

func InitPersistence(pool *pgxpool.Pool, log logger.Logger) Persistence {
	log.Info(context.Background(), "initializing payment persistence")

	repo := persistence.NewPaymentRepository(pool)
	log.Info(context.Background(), "payment persistence initialized")

	return Persistence{
		Payment: repo,
	}
}
