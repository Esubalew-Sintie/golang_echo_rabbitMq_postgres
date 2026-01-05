
---

## Payment Gateway API

A simple payment gateway service built with **Go**, **PostgreSQL**, and **RabbitMQ**.
It allows clients to create payments safely using **idempotency keys** and processes payments **asynchronously** in the background.

The service is designed to handle retries, avoid duplicate payments, and ensure reliable processing.

---

## How to Run the Project

### Prerequisites

* Go **1.21+**
* Docker & Docker Compose

---

### 1. Start Required Services

This starts PostgreSQL and RabbitMQ:

```bash
docker-compose -f scripts/docker-compose.yml up -d
```

Wait a few seconds for services to be ready:

```bash
sleep 10
```

---

### 2. Build and Run the Application

```bash
go build -o payment-gateway ./cmd
./payment-gateway
```

The API will be available at:

```
http://localhost:8080
```

---

### 3. Test the API

**Health Check**

```bash
curl http://localhost:8080/health
```

**Create a Payment**

```bash
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{"amount":100,"currency":"USD","idempotency_key":"test-001"}'
```

**Get Payment Status**

```bash
curl http://localhost:8080/api/v1/payments/{payment-id}
```

---
