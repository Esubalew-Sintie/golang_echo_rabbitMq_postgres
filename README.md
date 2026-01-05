💳 Payment Gateway API (Go)
A production-ready payment gateway built with Go, showcasing enterprise backend patterns such as idempotency, asynchronous processing, transaction safety, and reliable message delivery.

This project is ideal as:

A real-world backend reference
A system design / architecture showcase
A foundation for fintech or payment services
✨ Key Highlights
Idempotent payment creation (safe retries)
Asynchronous payment processing using RabbitMQ
ACID-compliant transactions with PostgreSQL
Row-level locking to prevent race conditions
Publisher confirms for guaranteed message delivery
Dockerized infrastructure
Clean Architecture & DDD-inspired structure
🧱 Tech Stack
Layer	Technology
Language	Go 1.21+
Web Framework	Echo
Database	PostgreSQL 15+
Messaging	RabbitMQ
ORM / SQL	sqlc
Containers	Docker & Docker Compose
🏗️ System Architecture
Client
  |
  v
HTTP API (Echo)
  |
  |-- PostgreSQL (Payments)
  |
  |-- RabbitMQ (Payment Events)
           |
           v
      Background Worker
           |
           v
      PostgreSQL (Status Update)
🚀 Quick Start
1️⃣ Start Infrastructure
docker-compose -f scripts/docker-compose.yml up -d
sleep 10
2️⃣ Build & Run
go build -o payment-gateway ./cmd
./payment-gateway
Service runs at:

http://localhost:8080
🔍 API Usage
Health Check
curl http://localhost:8080/health
Create Payment (Idempotent)
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.50,
    "currency": "USD",
    "idempotency_key": "order-123"
  }'
Response

{
  "status": 201,
  "message": "Payment created successfully",
  "data": {
    "payment_id": "uuid",
    "status": "PENDING"
  }
}
🔁 Sending the same request again with the same idempotency_key returns the same payment.

Get Payment Status
curl http://localhost:8080/api/v1/payments/{payment-id}
{
  "status": "SUCCESS",
  "amount": "100.50",
  "currency": "USD"
}
🔄 Payment Lifecycle
Client creates payment (PENDING)
Payment event published to RabbitMQ
Worker consumes message
Row locked using SELECT FOR UPDATE
Payment processed (success/failure)
Status updated atomically
✔️ Exactly-once processing ✔️ Safe under concurrency ✔️ Retry-friendly

🗄️ Database Schema
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('ETB', 'USD')),
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);
🐰 RabbitMQ Configuration
Exchange: payment.exchange
Queue: payment.processing
Routing Key: payment.processing
Delivery: Persistent + Publisher Confirms
📁 Project Structure
payment-gateway/
├── cmd/                # Application entrypoint
├── internal/
│   ├── handler/        # HTTP handlers
│   ├── service/        # Business logic
│   ├── domain/         # Entities
│   ├── storage/        # Database access
│   └── pkg/            # Shared utilities
├── database/
│   ├── migrations/
│   └── queries/
├── scripts/
│   ├── Dockerfile
│   └── docker-compose.yml
├── Makefile
└── README.md
🧪 Testing
go test ./...
go test -race ./...
go test -tags=integration ./...
Concurrent Idempotency Test
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/payments \
    -H "Content-Type: application/json" \
    -d '{"amount":50,"currency":"USD","idempotency_key":"load-test"}' &
done
wait
✔️ Only one payment is created

⚙️ Configuration
# Server
SERVER_ADDR=:8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=payment_gateway

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
🔐 Security Practices
Input validation
Parameterized SQL (sqlc)
Idempotency protection
No sensitive error leakage
Transaction isolation & locking
🚀 Production Ready Features
✔ Dockerized ✔ Horizontally scalable ✔ Fault-tolerant messaging ✔ Clean architecture ✔ Observability-friendly logging

🔮 Future Improvements
JWT authentication
Rate limiting
Webhooks
Prometheus metrics
OpenTelemetry tracing
Multi-currency expansion
📜 License
This project is intended for learning, demonstration, and production reference.

⭐ Why This Project Stands Out
This is not a toy API — it demonstrates how real payment systems are designed:

Safe retries
Asynchronous workflows
Strong consistency
Enterprise Go patterns
Built with Go, PostgreSQL, RabbitMQ & Docker 🚀 Clean. Scalable. Production-ready.