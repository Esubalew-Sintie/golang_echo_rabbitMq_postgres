.PHONY: help build run test clean docker-up docker-down docker-logs db-connect health-check api-test

# Default target
help: ## Show this help message
	@echo "Payment Gateway - Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# Build commands
build: ## Build the application
	go build -o payment-gateway ./cmd

run: ## Run the application (requires services to be running)
	go build -o payment-gateway ./cmd && ./payment-gateway

# Testing commands
test: ## Run unit tests
	go test ./...

test-integration: ## Run integration tests
	go test -tags=integration ./...

api-test: ## Run API health check
	curl -s http://localhost:8080/health | jq .

health-check: ## Check service health
	curl -s http://localhost:8080/health | jq .

# Docker commands
docker-up: ## Start all services (PostgreSQL + RabbitMQ)
	docker-compose -f scripts/docker-compose.yml up -d

docker-down: ## Stop all services
	docker-compose -f scripts/docker-compose.yml down

docker-logs: ## View all service logs
	docker-compose -f scripts/docker-compose.yml logs -f

docker-restart: ## Restart all services
	docker-compose -f scripts/docker-compose.yml restart

# Database commands
db-connect: ## Connect to PostgreSQL database
	docker exec -it payment-gateway-postgres psql -U postgres -d payment_gateway

db-status: ## Check database contents
	docker exec -it payment-gateway-postgres psql -U postgres -d payment_gateway -c "SELECT COUNT(*) as payments_count FROM payments;"

# Development commands
fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint (if installed)
	golangci-lint run

swagger: ## Generate Swagger documentation
	swag init -g ./cmd/main.go --parseInternal

clean: ## Clean build artifacts
	rm -f payment-gateway
	docker-compose -f scripts/docker-compose.yml down -v

# Combined commands
setup: docker-up ## Full setup: start services and build
	@echo "Waiting for services to be ready..."
	@sleep 10
	@make build

start: setup run ## Start everything: services + application

dev: ## Development mode with auto-restart (requires air)
	air

# Quick test commands
test-payment: ## Create a test payment
	curl -X POST http://localhost:8080/api/v1/payments \
		-H "Content-Type: application/json" \
		-d '{"amount": 100.50, "currency": "USD", "idempotency_key": "makefile-test"}'

test-idempotent: ## Test idempotency (run after test-payment)
	curl -X POST http://localhost:8080/api/v1/payments \
		-H "Content-Type: application/json" \
		-d '{"amount": 100.50, "currency": "USD", "idempotency_key": "makefile-test"}'