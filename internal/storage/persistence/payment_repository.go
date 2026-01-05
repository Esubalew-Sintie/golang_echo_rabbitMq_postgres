package persistence

import (
	"context"
	"database/sql"
	"errors"
	apperrors "payment-gateway/internal/constant/errors"
	"payment-gateway/internal/domain/entities"
	"payment-gateway/internal/storage"
	sqlc "payment-gateway/internal/storage/persistence/sqlc"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

type PaymentRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewPaymentRepository(db *pgxpool.Pool) storage.Persistence {
	return &PaymentRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *entities.Payment) error {
	params := sqlc.CreatePaymentParams{
		ID:             payment.ID,
		Amount:         payment.Amount,
		Currency:       payment.Currency,
		IdempotencyKey: payment.IdempotencyKey,
		Status:         string(payment.Status),
		CreatedAt:      pgtype.Timestamptz{Time: payment.CreatedAt, Valid: true},
	}

	_, err := r.queries.CreatePayment(ctx, params)
	if err != nil {
		if isUniqueViolation(err) {
			return apperrors.ErrDuplicateReference
		}
		return apperrors.ErrDatabaseOperationFailed
	}

	return nil
}

func (r *PaymentRepository) GetPaymentByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	payment, err := r.queries.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrPaymentNotFound
		}
		return nil, apperrors.ErrDatabaseOperationFailed
	}

	var processedAt *time.Time
	if payment.ProcessedAt.Valid {
		processedAt = &payment.ProcessedAt.Time
	}

	return &entities.Payment{
		ID:             payment.ID,
		Amount:         payment.Amount,
		Currency:       payment.Currency,
		IdempotencyKey: payment.IdempotencyKey,
		Status:         entities.PaymentStatus(payment.Status),
		CreatedAt:      payment.CreatedAt.Time,
		ProcessedAt:    processedAt,
	}, nil
}



func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status entities.PaymentStatus) error {
	params := sqlc.UpdatePaymentStatusParams{
		ID:          id,
		Status:      string(status),
		ProcessedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Status_2:    string(entities.PaymentStatusPending),
	}

	err := r.queries.UpdatePaymentStatus(ctx, params)
	if err != nil {
		return apperrors.ErrDatabaseOperationFailed
	}

	return nil
}

func (r *PaymentRepository) GetPaymentByIdempotencyKey(ctx context.Context, idempotencyKey string) (*entities.Payment, error) {
	payment, err := r.queries.GetPaymentByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrPaymentNotFound
		}
		return nil, apperrors.ErrDatabaseOperationFailed
	}

	var processedAt *time.Time
	if payment.ProcessedAt.Valid {
		processedAt = &payment.ProcessedAt.Time
	}

	return &entities.Payment{
		ID:             payment.ID,
		Amount:         payment.Amount,
		Currency:       payment.Currency,
		IdempotencyKey: payment.IdempotencyKey,
		Status:         entities.PaymentStatus(payment.Status),
		CreatedAt:      payment.CreatedAt.Time,
		ProcessedAt:    processedAt,
	}, nil
}

func (r *PaymentRepository) BeginTx(ctx context.Context) (storage.Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, apperrors.ErrDatabaseOperationFailed
	}
	return &pgxTransaction{tx: tx, ctx: ctx}, nil
}

type pgxTransaction struct {
	tx  pgx.Tx
	ctx context.Context
}

func (t *pgxTransaction) Commit() error {
	return t.tx.Commit(t.ctx)
}

func (t *pgxTransaction) Rollback() error {
	return t.tx.Rollback(t.ctx)
}

func (r *PaymentRepository) ProcessPaymentWithTransaction(ctx context.Context, paymentID uuid.UUID, newStatus entities.PaymentStatus) error {

	payment, err := r.GetPaymentByID(ctx, paymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrPaymentNotFound
		}
		return apperrors.ErrDatabaseOperationFailed
	}

	paymentEntity := entities.Payment{
		ID:     payment.ID,
		Status: entities.PaymentStatus(payment.Status),
	}

	if paymentEntity.IsTerminalState() {
		return apperrors.ErrPaymentAlreadyProcessed
	}

	if !paymentEntity.CanBeProcessed() {
		return apperrors.ErrInvalidPaymentStatus
	}

	updateParams := sqlc.UpdatePaymentStatusParams{
		ID:          paymentID,
		Status:      string(newStatus),
		ProcessedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Status_2:    string(entities.PaymentStatusPending),
	}

	err = r.queries.UpdatePaymentStatus(ctx, updateParams)
	if err != nil {
		return apperrors.ErrDatabaseOperationFailed
	}

	return nil
}
