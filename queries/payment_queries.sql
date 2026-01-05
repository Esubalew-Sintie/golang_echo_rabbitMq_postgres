-- name: CreatePayment :one
INSERT INTO payments (id, amount, currency, reference, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, amount, currency, reference, status, created_at, processed_at, updated_at;

-- name: GetPaymentByID :one
SELECT id, amount, currency, reference, status, created_at, processed_at, updated_at
FROM payments
WHERE id = $1;

-- name: GetPaymentByIDForUpdate :one
SELECT id, amount, currency, reference, status, created_at, processed_at, updated_at
FROM payments
WHERE id = $1
FOR UPDATE;

-- name: UpdatePaymentStatus :exec
UPDATE payments
SET status = $2, processed_at = $3, updated_at = NOW()
WHERE id = $1 AND status = $4;

-- name: GetPaymentByReference :one
SELECT id, amount, currency, reference, status, created_at, processed_at, updated_at
FROM payments
WHERE reference = $1;

-- name: ListPayments :many
SELECT id, amount, currency, reference, status, created_at, processed_at, updated_at
FROM payments
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPaymentsByStatus :many
SELECT id, amount, currency, reference, status, created_at, processed_at, updated_at
FROM payments
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPayments :one
SELECT COUNT(*) as total
FROM payments;

-- name: CountPaymentsByStatus :one
SELECT COUNT(*) as total
FROM payments
WHERE status = $1;