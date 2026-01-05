# Payment Gateway API

A production-ready payment gateway service built with Go that demonstrates enterprise-grade architecture patterns including idempotent operations, asynchronous processing, and reliable message delivery.

## 🚀 Quick Commands

```bash
# Start everything (PostgreSQL + RabbitMQ + App)
docker-compose -f scripts/docker-compose.yml up -d && go run ./cmd/main.go

# Test API health
curl http://localhost:8080/health

# Create payment
curl -X POST http://localhost:8080/api/v1/payments -H "Content-Type: application/json" -d '{"amount": 100.50, "currency": "USD", "idempotency_key": "test-123"}'

# Get payment status
curl http://localhost:8080/api/v1/payments/{payment-id}
```

## 🚀 Features

- ✅ **Payment Creation API** - Create payments with validation
- ✅ **Idempotent Operations** - Safe retry handling with unique idempotency keys
- ✅ **Asynchronous Processing** - Background payment processing via RabbitMQ
- ✅ **Database Transactions** - ACID compliance with PostgreSQL
- ✅ **Row-Level Locking** - Prevents race conditions during processing
- ✅ **Publisher Confirms** - Guaranteed message delivery
- ✅ **Docker Containerization** - Complete containerized setup
- ✅ **Comprehensive Error Handling** - Detailed API responses

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP API      │    │   RabbitMQ      │    │   Worker        │
│   (Echo)        │◄──►│   (Message      │◄──►│   (Go Routine)  │
│                 │    │    Queue)       │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                        │
         ▼                        ▼
┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │   PostgreSQL    │
│   (Payments)    │    │   (Payments)    │
└─────────────────┘    └─────────────────┘
```

## 📋 Prerequisites

- **Go 1.21+**
- **PostgreSQL 15+**
- **RabbitMQ 3.8+**
- **Docker & Docker Compose**

## 🛠️ Quick Start

### One-Line Setup (Everything)

```bash
# Complete setup: clone, start services, build and run
git clone <repository-url> && cd golang_payment/payment-gateway && \
docker-compose -f scripts/docker-compose.yml up -d && \
sleep 10 && go build -o payment-gateway ./cmd && ./payment-gateway
```

### Detailed Setup

```bash
git clone <repository-url>
cd golang_payment/payment-gateway
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL and RabbitMQ
docker-compose -f scripts/docker-compose.yml up -d

# Verify services are running
docker-compose -f scripts/docker-compose.yml ps
```

### 3. Run Database Migrations

```bash
# The migrations run automatically via Docker entrypoint
# Or run manually if needed:
docker exec -it payment-gateway-postgres psql -U postgres -d payment_gateway -f /docker-entrypoint-initdb.d/000001_create_payments_table.up.sql
```

### 4. Build and Run

```bash
# Build the application
go build -o payment-gateway ./cmd

# Run the service
./payment-gateway
```

The API will be available at `http://localhost:8080`

### Manual API Testing

Test the key features manually:

```bash
# 1. Health check
curl http://localhost:8080/health

# 2. Create payment
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.50, "currency": "USD", "idempotency_key": "test-001"}'

# 3. Test idempotency (same request again)
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.50, "currency": "USD", "idempotency_key": "test-001"}'

# 4. Get payment status (wait a few seconds for processing)
curl http://localhost:8080/api/v1/payments/{payment-id-from-response}

# 5. Test validation (invalid currency)
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.00, "currency": "INVALID", "idempotency_key": "test-002"}'
```

## 🧪 Manual Testing the API

### Create a Payment

```bash
# First payment creation
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.50,
    "currency": "USD",
    "idempotency_key": "test-payment-001"
  }'

# Response: 201 Created
{
  "status": 201,
  "message": "Payment created successfully",
  "data": {
    "payment_id": "uuid-here",
    "status": "PENDING",
    "created_at": "2024-01-05T10:30:00Z"
  }
}
```

### Test Idempotency (Duplicate Request)

```bash
# Send the same request again - should return existing payment
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.50,
    "currency": "USD",
    "idempotency_key": "test-payment-001"
  }'

# Response: 200 OK (existing payment returned)
{
  "status": 200,
  "message": "Payment already exists",
  "data": {
    "payment_id": "same-uuid-as-above",
    "status": "PENDING",
    "created_at": "2024-01-05T10:30:00Z"
  }
}
```

### Get Payment Status

```bash
# Replace with actual payment ID
curl -X GET http://localhost:8080/api/v1/payments/{payment-id}

# Response
{
  "status": 200,
  "message": "Payment retrieved successfully",
  "data": {
    "id": "uuid-here",
    "amount": "100.50",
    "currency": "USD",
    "idempotency_key": "test-payment-001",
    "status": "SUCCESS",  // or "FAILED" after processing
    "created_at": "2024-01-05T10:30:00Z",
    "processed_at": "2024-01-05T10:30:05Z"
  }
}
```

### Health Check

```bash
curl -X GET http://localhost:8080/health

# Response
{
  "status": 200,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "service": "payment-gateway",
    "version": "1.0.0"
  }
}
```

## 📊 Database Schema

```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('ETB', 'USD')),
    idempotency_key VARCHAR(255) NOT NULL UNIQUE,  -- Idempotency key
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'SUCCESS', 'FAILED')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_payments_idempotency_key ON payments(idempotency_key);
CREATE INDEX idx_payments_status ON payments(status);
```

## 🔄 Payment Processing Flow

1. **Payment Creation**
   - Client sends POST `/api/v1/payments` with `idempotency_key`
   - Database unique constraint ensures no duplicates
   - Payment status set to `PENDING`
   - Message published to RabbitMQ

2. **Asynchronous Processing**
   - Worker consumes message from RabbitMQ
   - Uses row-level locking (`SELECT ... FOR UPDATE`) to prevent race conditions
   - Simulates payment processing (90% success rate)
   - Updates payment status to `SUCCESS` or `FAILED`

3. **Idempotency Guarantee**
   - Same `idempotency_key` always returns same payment
   - Database constraint prevents duplicate payments
   - Safe for client retries and network failures

## 🐰 RabbitMQ Setup

The application uses RabbitMQ for reliable message delivery:

- **Exchange**: `payment.exchange`
- **Queue**: `payment.processing`
- **Routing Key**: `payment.processing`
- **Publisher Confirms**: Enabled for guaranteed delivery
- **Persistent Messages**: Messages survive broker restarts

## 🧪 Advanced Testing Scenarios

### Load Testing Idempotency

```bash
# Test concurrent requests with same idempotency key
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/payments \
    -H "Content-Type: application/json" \
    -d '{"amount": 50.00, "currency": "USD", "idempotency_key": "load-test-001"}' &
done
wait

# Result: Only 1 payment created, others return the same payment
```

### Message Delivery Testing

```bash
# Monitor RabbitMQ management UI
open http://localhost:15672 (guest/guest)

# Check message delivery and consumption
# Messages should be processed exactly once
```

### Database Transaction Testing

```bash
# Simulate network failure during payment creation
# Restart the service mid-request
# Verify no orphaned payments or messages
```

## 📋 Command Reference

### Core Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make start` | Start everything (services + app) |
| `make api-test` | Check API health |
| `make docker-up` | Start PostgreSQL + RabbitMQ |
| `make build` | Build the application |
| `make run` | Run the application |
| `make clean` | Clean up everything |

#### Quick Start with Make

```bash
# Start everything in one command
make start

# Run tests
make api-test

# Check what's available
make help
```

### Development Commands

```bash
# Build application
go build -o payment-gateway ./cmd

# Run application
./payment-gateway

# Run tests
go test ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...

# Run linter
golangci-lint run

# Regenerate SQLC code
sqlc generate
```

### Docker Commands

```bash
# Start services
docker-compose -f scripts/docker-compose.yml up -d

# Stop services
docker-compose -f scripts/docker-compose.yml down

# View logs
docker-compose -f scripts/docker-compose.yml logs -f payment-gateway
docker-compose -f scripts/docker-compose.yml logs -f postgres
docker-compose -f scripts/docker-compose.yml logs -f rabbitmq

# Restart specific service
docker-compose -f scripts/docker-compose.yml restart payment-gateway

# Scale application
docker-compose -f scripts/docker-compose.yml up -d --scale payment-gateway=3
```

### Database Commands

```bash
# Connect to database
docker exec -it payment-gateway-postgres psql -U postgres -d payment_gateway

# Check payments
SELECT * FROM payments;

# Check by idempotency key
SELECT * FROM payments WHERE idempotency_key = 'test-123';

# View indexes
\di
```

## 🔧 Configuration

### Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=payment_gateway
DB_SSLMODE=disable

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_VHOST=/

# Application
SERVER_ADDR=:8080
```

### Docker Environment

The `docker-compose.yml` provides all necessary environment variables for containerized deployment.

## 📈 Monitoring & Observability

### Health Checks

- **Service Health**: `GET /health`
- **Database Connectivity**: Built-in connection pooling
- **Message Queue**: Publisher confirms ensure delivery

### Logs

The application provides structured logging for:

- Payment creation attempts
- Idempotency key collisions
- Message publishing confirmations
- Processing worker activity
- Database transaction outcomes

## 🏗️ Project Structure

```
payment-gateway/
├── cmd/                          # Application entrypoint
├── internal/
│   ├── constant/
│   │   ├── errors/              # Error definitions
│   │   └── model/
│   │       ├── dto/             # Data transfer objects
│   │       └── response/        # API response structures
│   ├── domain/
│   │   └── entities/            # Domain models
│   ├── handler/                 # HTTP handlers
│   ├── service/                 # Business logic
│   ├── storage/                 # Data persistence
│   └── pkg/                     # Shared utilities
├── database/
│   ├── migrations/              # Database schema
│   └── queries/                 # SQL queries
├── scripts/
│   ├── Dockerfile               # Container definition
│   └── docker-compose.yml       # Infrastructure setup
├── Makefile                     # Build and development commands
└── README.md                    # This file
```

## 🔒 Security Considerations

- **Input Validation**: Comprehensive validation using ozzo-validation
- **SQL Injection Protection**: Parameterized queries via sqlc
- **Idempotency Keys**: Prevent duplicate operations
- **Transaction Safety**: ACID compliance prevents inconsistencies
- **Error Handling**: No sensitive information leaked in responses

## 🚀 Production Deployment

### Docker Deployment

```bash
# Build and deploy
docker-compose -f scripts/docker-compose.yml up -d

# Scale the worker service
docker-compose -f scripts/docker-compose.yml up -d --scale payment-gateway=3
```

### Health Monitoring

```bash
# Check service health
curl http://localhost:8080/health

# Monitor Docker containers
docker-compose -f scripts/docker-compose.yml ps

# View application logs
docker-compose -f scripts/docker-compose.yml logs -f payment-gateway

# View database logs
docker-compose -f scripts/docker-compose.yml logs -f postgres

# View RabbitMQ logs
docker-compose -f scripts/docker-compose.yml logs -f rabbitmq
```

### Troubleshooting

#### Common Issues

**Port 8080 already in use:**
```bash
# Kill process using port 8080
lsof -ti:8080 | xargs kill -9

# Or run on different port
SERVER_ADDR=:8081 ./payment-gateway
```

**Database connection failed:**
```bash
# Check if PostgreSQL is running
docker-compose -f scripts/docker-compose.yml ps postgres

# Restart database
docker-compose -f scripts/docker-compose.yml restart postgres
```

**RabbitMQ connection failed:**
```bash
# Check RabbitMQ status
docker-compose -f scripts/docker-compose.yml ps rabbitmq

# Access RabbitMQ management UI
open http://localhost:15672 (guest/guest)
```


#### Database Issues

**Reset database:**
```bash
# Stop and remove containers
docker-compose -f scripts/docker-compose.yml down -v

# Restart fresh
docker-compose -f scripts/docker-compose.yml up -d
```

**Manual database inspection:**
```bash
# Connect to database
docker exec -it payment-gateway-postgres psql -U postgres -d payment_gateway

# Check payments table
SELECT * FROM payments;

# Check indexes
\di
```

## 🤝 API Reference

### POST `/api/v1/payments`
Create a new payment with idempotency guarantee.

**Request Body:**
```json
{
  "amount": 100.50,
  "currency": "USD",
  "idempotency_key": "unique-reference-123"
}
```

**Response (201 Created):**
```json
{
  "status": 201,
  "message": "Payment created successfully",
  "data": {
    "payment_id": "uuid-here",
    "status": "PENDING",
    "created_at": "2024-01-05T10:30:00Z"
  }
}
```

### GET `/api/v1/payments/{id}`
Retrieve payment details by ID.

**Response (200 OK):**
```json
{
  "status": 200,
  "message": "Payment retrieved successfully",
  "data": {
    "id": "uuid-here",
    "amount": "100.50",
    "currency": "USD",
    "idempotency_key": "unique-reference-123",
    "status": "SUCCESS",
    "created_at": "2024-01-05T10:30:00Z",
    "processed_at": "2024-01-05T10:30:05Z"
  }
}
```

### GET `/health`
Service health check endpoint.

## 📝 Development

### Adding New Features

1. **Database Changes**: Add migrations in `database/migrations/`
2. **SQL Queries**: Add queries in `database/queries/`
3. **Regenerate sqlc**: `sqlc generate`
4. **Implement Service**: Add business logic in `internal/service/`
5. **Add Handler**: Create HTTP endpoints in `internal/handler/`

### Running Tests

```bash
# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./...

# Run API tests
# Manual API testing commands above

# Development workflow (with hot reload)
go install github.com/cosmtrek/air@latest
air
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run linter
golangci-lint run
```

## 🔮 Future Enhancements

- **JWT Authentication**: Secure API endpoints
- **Rate Limiting**: Prevent abuse
- **Metrics & Monitoring**: Prometheus integration
- **Distributed Tracing**: OpenTelemetry support
- **API Versioning**: Support multiple API versions
- **Webhook Notifications**: Real-time payment status updates
- **Admin Dashboard**: Payment management interface
- **Multi-Currency Support**: Extended currency handling
- **Payment Methods**: Credit cards, digital wallets, etc.

## 🤝 Contributing

This implementation demonstrates production-ready patterns for:

- **Clean Architecture**: Separation of concerns
- **Domain-Driven Design**: Business logic encapsulation
- **CQRS Pattern**: Command and query separation
- **Event-Driven Architecture**: Asynchronous processing
- **Transactional Sagas**: Distributed transaction management
- **Idempotent Operations**: Safe retry handling

## 📞 Support

For questions about this payment gateway implementation:

- Review the code comments for detailed explanations
- Check the test script for usage examples
- Examine the database schema for data relationships
- Monitor logs for troubleshooting information

---

**🎉 Happy coding! This payment gateway showcases enterprise-grade Go development practices.**

## 📄 License

This project demonstrates enterprise-grade Go development patterns for payment processing systems.

---

**Built with Go, PostgreSQL, RabbitMQ, and Docker for reliable, scalable payment processing!** 🚀