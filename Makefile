.PHONY: build build-gateway build-service build-booking run run-gateway run-service run-booking run-all run-both test clean docker-build docker-run docker-build-booking docker-run-booking deps migrate-up migrate-down seed-db seed-reset health-check rabbitmq-status redis-status status dev-run

# Build variables
BINARY_NAME_GATEWAY=api-gateway
BINARY_NAME_SERVICE=api-service
BINARY_NAME_BOOKING=booking-service
DOCKER_IMAGE=airline-booking:latest

# Default target
help:
	@echo "Available targets:"
	@echo ""
	@echo "🔨 Build:"
	@echo "  build              - Build all services"
	@echo "  build-gateway      - Build API gateway"
	@echo "  build-service      - Build API service"
	@echo "  build-booking      - Build booking service"
	@echo ""
	@echo "🚀 Run:"
	@echo "  run                - Run all services (alias for run-all)"
	@echo "  run-all            - Run all services (gateway, api, booking)"
	@echo "  run-both           - Run API service and gateway only"
	@echo "  run-gateway        - Run API gateway only"
	@echo "  run-service        - Run API service only" 
	@echo "  run-booking        - Run booking service only"
	@echo "  dev-run            - Run all services in development mode"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-build-booking - Build booking service Docker image"
	@echo "  docker-run-booking - Run booking service Docker container"
	@echo ""
	@echo "💾 Database:"
	@echo "  migrate-up         - Run database migrations"
	@echo "  migrate-down       - Rollback database migrations"
	@echo "  seed-db            - Seed database with demo data"
	@echo "  seed-reset         - Reset database and reseed"
	@echo "  db-status          - Check database status"
	@echo ""
	@echo "🏥 Health & Status:"
	@echo "  health-check       - Check all service health endpoints"
	@echo "  rabbitmq-status    - Check RabbitMQ queue status"
	@echo "  redis-status       - Check Redis cache status"
	@echo "  status             - Full system status check"
	@echo ""
	@echo "🔧 Development:"
	@echo "  dev-setup          - Setup complete development environment"
	@echo "  demo-setup         - Setup demo environment with sample data"
	@echo "  dev-up             - Start Docker services"
	@echo "  dev-down           - Stop Docker services"
	@echo ""
	@echo "🧹 Maintenance:"
	@echo "  test               - Run all tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  clean              - Clean build artifacts"
	@echo "  deps               - Install dependencies"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linter"
	@echo "  vet                - Run go vet"

# Build the application
build: build-gateway build-service build-booking

build-gateway:
	go build -o bin/$(BINARY_NAME_GATEWAY) ./cmd/api-gateway

build-service:
	go build -o bin/$(BINARY_NAME_SERVICE) ./cmd/api-service

build-booking:
	go build -o bin/$(BINARY_NAME_BOOKING) ./cmd/booking-service

# Run the application
run: run-all

run-gateway:
	go run ./cmd/api-gateway

run-service:
	go run ./cmd/api-service

run-booking:
	go run ./cmd/booking-service

run-both:
	@echo "Starting API Service on port 8080..."
	@go run ./cmd/api-service &
	@sleep 2
	@echo "Starting API Gateway on port 8000..."
	@go run ./cmd/api-gateway

run-all:
	@echo "Starting API Service on port 8080..."
	@go run ./cmd/api-service &
	@echo "Starting Booking Service on port 8081..."
	@go run ./cmd/booking-service &
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
	rm -f bin/$(BINARY_NAME_BOOKING)
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

# Health checks
health-check:
	@echo "Checking service health..."
	@echo "API Gateway (port 8000):"
	@curl -s http://localhost:8000/health || echo "  ❌ Not responding"
	@echo ""
	@echo "API Service (port 8080):"
	@curl -s http://localhost:8080/health || echo "  ❌ Not responding"
	@echo ""
	@echo "Booking Service (port 8081):"
	@curl -s http://localhost:8081/health || echo "  ❌ Not responding"

# RabbitMQ status
rabbitmq-status:
	@echo "Checking RabbitMQ queues..."
	@curl -s -u admin:admin123 http://localhost:15672/api/queues/%2f/booking_flight_queue | jq '.messages' 2>/dev/null || echo "RabbitMQ not accessible"

# Redis status  
redis-status:
	@echo "Checking Redis booking cache..."
	@redis-cli -h localhost -p 6379 keys "booking_status:*" 2>/dev/null || echo "Redis not accessible"

# Full status check
status: health-check rabbitmq-status redis-status db-status

# Docker build booking service
docker-build-booking:
	docker build --target booking-service -t airline-booking-booking:latest .

# Run booking service in Docker
docker-run-booking:
	docker run -p 8081:8081 airline-booking-booking:latest

# Development mode with all services
dev-run:
	@echo "Starting all services in development mode..."
	@echo "Make sure to run 'make dev-setup' first if this is your first time!"
	@echo ""
	make run-all