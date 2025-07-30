# ✈️ Airline Booking System

A high-performance, microservices-based airline booking system built with Go, featuring intelligent seat assignment, optimistic locking for data consistency, and asynchronous booking processing.

## 🏗️ Architecture

### Microservices
- **API Gateway** (Port 8000) - Request routing and load balancing
- **API Service** (Port 8080) - REST API for flight search and booking requests
- **Booking Service** (Port 8081) - Asynchronous booking processing with seat assignment

### Infrastructure
- **PostgreSQL** - Primary database with read replicas
- **Redis** - Caching layer and session management
- **RabbitMQ** - Message queuing for asynchronous processing
- **Docker** - Containerized deployment

## 🚀 Quick Start

### 1. View All Available Commands
```bash
make help
# or just
make
```

### 2. Build All Services
```bash
make build
```

### 3. Run All Services (including new booking service)
```bash
# Start infrastructure first
make dev-up

# Wait for services to be ready, then
make run-all
```

### 4. Complete Setup (First Time)
```bash
# Setup entire development environment
make demo-setup
```

This will:
- Start all Docker services (PostgreSQL, Redis, RabbitMQ)
- Run database migrations
- Seed database with sample data
- Set up demo user accounts

## 📋 Prerequisites

- **Go 1.21+**
- **Docker & Docker Compose**
- **Make**
- **PostgreSQL** (for local development)
- **Redis** (for local development)
- **RabbitMQ** (for local development)

## 🛠️ Development

### Environment Setup
```bash
# Clone the repository
git clone <repository-url>
cd airline_booking

# Setup development environment
make dev-setup

# Build all services
make build

# Run all services
make run-all
```

### Individual Service Management
```bash
# Build individual services
make build-gateway      # API Gateway
make build-service      # API Service  
make build-booking      # Booking Service

# Run individual services
make run-gateway        # Port 8000
make run-service        # Port 8080
make run-booking        # Port 8081
```

### Health Monitoring
```bash
make health-check       # Check all service endpoints
make status            # Complete system status
make rabbitmq-status   # Queue status
make redis-status      # Cache status
make db-status         # Database status
```

## 🏛️ Database Schema

### Core Entities
- **Users** - Customer accounts and authentication
- **Airlines** - Airline companies and metadata
- **Aircraft** - Aircraft models and configurations
- **Airports** - Airport codes and information
- **Flights** - Flight definitions and routes
- **Flight Schedules** - Specific flight instances
- **Flight Inventory** - Seat availability by class
- **Seat Maps** - Aircraft seat configurations
- **Seat Assignments** - Individual seat reservations
- **Bookings** - Customer reservations
- **Booking Passengers** - Passenger details per booking

### Key Features
- **Optimistic Locking** - Version-controlled updates prevent race conditions
- **Read Replicas** - Separate read/write database connections
- **Seed Data** - Sample airlines, airports, flights, and users

## 🔄 Booking Flow

### 1. Flight Search
```
Client → API Gateway → API Service → Database (Read Replica) → Redis Cache
```

### 2. Booking Request
```
Client → API Gateway → API Service → RabbitMQ Queue
```

### 3. Booking Processing
```
RabbitMQ → Booking Service → Seat Assignment → Database (Master) → Redis Cache
```

### 4. Status Updates
```
Booking Service → Redis Cache → Client (via status API)
```

## 🔐 Optimistic Locking

### Implementation
The system uses optimistic locking to prevent race conditions during high-concurrency operations:

- **Seat Assignments** - Version-controlled seat reservations
- **Flight Inventory** - Atomic seat count updates
- **Cache Invalidation** - Automatic cache clearing on version changes

### Conflict Resolution
- **Automatic Retries** - Up to 3 attempts with exponential backoff
- **Intelligent Fallback** - Alternative seat selection on conflicts
- **Graceful Degradation** - Clear error messages for failed bookings

## 📊 API Endpoints

### Authentication
```bash
POST /api/login
```

### Flight Search
```bash
POST /api/flights/search
```

### Booking Management
```bash
POST /api/flights/bookings        # Create booking
GET  /api/flights/bookings/:uuid  # Get booking status
```

### Health Checks
```bash
GET /health                       # Service health
```

## 🔧 Configuration

### Environment Variables
```bash
# Database
DATABASE_DSN="host=localhost user=postgres password=postgres dbname=airline_booking port=5432 sslmode=disable"
SLAVE_DATABASE_DSN="..."

# Redis
REDIS_HOST="localhost"
REDIS_PORT="6379"

# RabbitMQ
RABBITMQ_URL="amqp://admin:admin123@localhost:5672/"

# Booking Service Retry Configuration
BOOKING_MAX_RETRIES=3
BOOKING_INITIAL_DELAY_MS=100
BOOKING_MAX_DELAY_MS=2000
BOOKING_BACKOFF_MULTIPLIER=2.0

# Security
JWT_SECRET="your-secret-key"
```

### Docker Compose Services
- **postgres-master** (5432) - Primary database
- **redis** (6379) - Cache server
- **rabbitmq** (5672, 15672) - Message broker with management UI

## 🧪 Testing

### Run Tests
```bash
make test              # All tests
make test-coverage     # With coverage report
```

### API Testing
Use the Bruno collection in `docs/bruno/airline booking/`:
- User authentication
- Flight search
- Booking creation
- Status checking

### Demo User Accounts
```
admin@airline.com / password123 (Admin)
john.doe@example.com / password123 (Regular User)
jane.smith@example.com / password123 (Regular User)
demo@example.com / password123 (Demo User)
```

## 📈 Monitoring & Observability

### Service Health
- **Health Endpoints** - `/health` on each service
- **Metrics Collection** - Retry attempts, success rates, conflict types
- **Structured Logging** - JSON logs with correlation IDs

### Queue Monitoring
- **RabbitMQ Management** - http://localhost:15672 (admin/admin123)
- **Queue Metrics** - Message rates, consumer lag
- **Dead Letter Queues** - Failed message handling

### Cache Monitoring
- **Redis CLI** - Monitor cache hit rates and invalidations
- **Cache Keys** - `flight_search:*`, `booking_status:*`, `inventory_version:*`

## 🚦 Performance Features

### Caching Strategy
- **Flight Search Results** - 5-minute TTL with version-based invalidation
- **Booking Status** - 24-hour TTL with real-time updates
- **Inventory Versions** - Persistent version tracking

### Concurrency Control
- **Connection Pooling** - Optimized database connections
- **Message Prefetch** - Controlled RabbitMQ concurrency
- **Optimistic Locking** - Non-blocking read operations

### Scalability
- **Horizontal Scaling** - Multiple booking service instances
- **Load Balancing** - API Gateway request distribution
- **Read Replicas** - Separate read/write database access

## 🔄 Retry Logic

### Intelligent Retries
The booking service implements sophisticated retry logic:

```go
// Configurable retry parameters
MaxRetries: 3
InitialDelay: 100ms
MaxDelay: 2000ms
BackoffMultiplier: 2.0
```

### Retry Scenarios
- **Version Conflicts** - Optimistic locking failures
- **Seat Unavailability** - Real-time seat competition
- **Database Connectivity** - Transient connection issues
- **Cache Failures** - Redis connectivity problems

### Fallback Strategies
- **Seat Selection** - Alternative seat preferences
- **Preference Relaxation** - Less strict requirements on retries
- **Randomization** - Jittered seat selection to avoid conflicts

## 📝 Development Workflow

### Daily Development
```bash
# Start infrastructure
make dev-up

# Build and run services
make build
make run-all

# Check system health
make status

# Run tests
make test
```

### Code Quality
```bash
make fmt               # Format code
make lint              # Run linter
make vet               # Run go vet
```

### Database Management
```bash
make migrate-up        # Apply migrations
make migrate-down      # Rollback migrations
make seed-db           # Load sample data
make seed-reset        # Reset and reseed
```

## 🐳 Docker Deployment

### Build Images
```bash
make docker-build              # API Gateway/Service
make docker-build-booking      # Booking Service
```

### Run Containers
```bash
make docker-run                # API services
make docker-run-booking        # Booking service
```

### Full Stack
```bash
docker-compose up -d           # All services
```

## 📚 Documentation

### Additional Documentation
- **[Optimistic Locking](OPTIMISTIC_LOCKING.md)** - Detailed concurrency control implementation
- **[Makefile Usage](MAKEFILE_USAGE.md)** - Complete command reference
- **[Test Plan](test_booking_flow.md)** - End-to-end testing guide

### API Documentation
- **Bruno Collections** - `docs/bruno/airline booking/`
- **Database Schema** - `db.dbml`
- **Sequence Diagrams** - `*.mermaid` files

## 🤝 Contributing

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Run `make dev-setup` for complete environment
4. Make changes with tests
5. Run `make test && make lint`
6. Submit pull request

### Code Standards
- **Go fmt** - Standard formatting
- **golangci-lint** - Comprehensive linting
- **Test Coverage** - Minimum 80% coverage
- **Documentation** - Update relevant docs

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

#### Services won't start
```bash
make clean             # Clean build artifacts
make deps              # Reinstall dependencies
make build             # Rebuild everything
```

#### Database connection issues
```bash
make db-status         # Check database connectivity
make dev-down && make dev-up  # Restart infrastructure
```

#### Queue processing problems
```bash
make rabbitmq-status   # Check queue status
# Visit http://localhost:15672 for RabbitMQ management
```

#### Cache inconsistencies
```bash
make redis-status      # Check Redis connectivity
redis-cli flushall     # Clear all cache (development only)
```

### Getting Help

- **Health Checks** - `make status` for complete system overview
- **Logs** - Check service logs for detailed error information
- **Monitoring** - Use RabbitMQ and Redis management interfaces
- **Documentation** - Refer to additional docs in the repository

---

## 🎯 Key Features Summary

✅ **Microservices Architecture** - Scalable, maintainable service design  
✅ **Optimistic Locking** - Race condition prevention with version control  
✅ **Intelligent Retries** - Sophisticated conflict resolution  
✅ **Cache Strategy** - Multi-layer caching with automatic invalidation  
✅ **Asynchronous Processing** - Queue-based booking workflow  
✅ **Health Monitoring** - Comprehensive service observability  
✅ **Docker Support** - Containerized deployment ready  
✅ **Developer Experience** - Rich tooling and documentation  

**Built for high-concurrency airline booking scenarios with enterprise-grade reliability! 🚀**