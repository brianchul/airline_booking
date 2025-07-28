.PHONY: build build-gateway build-service run run-gateway run-service run-both test clean docker-build docker-run deps migrate-up migrate-down seed-db seed-reset

# Build variables
BINARY_NAME_GATEWAY=api-gateway
BINARY_NAME_SERVICE=api-service
DOCKER_IMAGE=airline-booking:latest

# Build the application
build: build-gateway build-service

build-gateway:
	go build -o bin/$(BINARY_NAME_GATEWAY) ./cmd/api-gateway

build-service:
	go build -o bin/$(BINARY_NAME_SERVICE) ./cmd/api-service

# Run the application
run: run-both

run-gateway:
	go run ./cmd/api-gateway

run-service:
	go run ./cmd/api-service

run-both:
	@echo "Starting API Service on port 8080..."
	@go run ./cmd/api-service &
	@sleep 2
	@echo "Starting API Gateway on port 8000..."
	@go run ./cmd/api-gateway

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	go clean
	rm -f bin/$(BINARY_NAME_GATEWAY)
	rm -f bin/$(BINARY_NAME_SERVICE)
	rm -f coverage.out

# Install dependencies
deps:
	go mod download
	go mod tidy

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Docker run
docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE)

# Start all services
dev-up:
	docker-compose up -d

# Stop all services
dev-down:
	docker-compose down

# Database migrations
migrate-up:
	migrate -path migrations -database "postgresql://postgres:rootpass@localhost:5432/airline_booking?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:rootpass@localhost:5432/airline_booking?sslmode=disable" down

# Lint code
lint:
	golangci-lint run

# Generate code
generate:
	go generate ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Database seeding
seed-db:
	@echo "Seeding database with demo data..."
	./seeds/run_seeds.sh

seed-reset:
	@echo "Resetting database and seeding with fresh demo data..."
	make migrate-down
	make migrate-up
	make seed-db

# Development setup (complete environment)
dev-setup:
	@echo "Setting up development environment..."
	make dev-up
	@echo "Waiting for database to be ready..."
	sleep 10
	make migrate-up
	make seed-db
	@echo "Development environment ready!"

# Quick demo setup
demo-setup:
	@echo "Setting up demo environment with sample data..."
	make dev-setup
	@echo "Demo environment ready!"
	@echo ""
	@echo "Demo user accounts (password: password123):"
	@echo "  - admin / admin@airline.com"
	@echo "  - john.doe / john.doe@example.com" 
	@echo "  - jane.smith / jane.smith@example.com"
	@echo "  - demo.user / demo@example.com"
	@echo ""
	@echo "Services available:"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379" 
	@echo "  - RabbitMQ Management: http://localhost:15672 (admin/admin123)"
	@echo "  - Elasticsearch: http://localhost:9200"
	@echo "  - Grafana: http://localhost:3000 (admin/admin)"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Jaeger: http://localhost:16686"

# Database status check
db-status:
	@echo "Checking database status..."
	@psql -h localhost -p 5432 -U postgres -d airline_booking -c "SELECT 'Users: ' || COUNT(*) FROM users; SELECT 'Airlines: ' || COUNT(*) FROM airlines; SELECT 'Flights: ' || COUNT(*) FROM flights; SELECT 'Bookings: ' || COUNT(*) FROM bookings;" 2>/dev/null || echo "Database not accessible"