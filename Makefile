.PHONY: build run test clean docker-build docker-run migrate-up migrate-down sqlc-gen sqlc-clean mocks

# Build the application
build:
	go build -o bin/payment-gateway ./cmd

# Run the application
run:
	go run ./cmd

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -f scripts/Dockerfile -t payment-gateway .

# Run with Docker Compose
docker-run:
	docker-compose -f scripts/docker-compose.yml up --build

# Run database migrations (placeholder)
migrate-up:
	@echo "Running database migrations..."
	# TODO: Implement migration command

# Rollback database migrations (placeholder)
migrate-down:
	@echo "Rolling back database migrations..."
	# TODO: Implement migration rollback command

# Generate database code with sqlc
sqlc-gen:
	@echo "Generating database code with sqlc..."
	sqlc generate

# Clean generated sqlc code
sqlc-clean:
	@echo "Cleaning generated sqlc code..."
	rm -f internal/storage/persistence/db.go internal/storage/persistence/models.go internal/storage/persistence/payment_queries.sql.go

# Generate mocks (if using mockery)
mocks:
	@echo "Generating mocks..."
	# TODO: Add mock generation commands
