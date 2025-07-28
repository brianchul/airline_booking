# Airline Booking System - Folder Structure

This document describes the purpose and contents of each directory in the airline booking system.

## Project Root

```
airline-booking/
├── cmd/                    # Application entry points
│   ├── api-gateway/        # API reverse proxy service with jwt validation
|   ├── api-service/        # API services to connect remain service
│   ├── booking-service/    # Booking service
│   ├── flight-service/     # Flight service
│   ├── notification-service/ # Notification service
│   └── payment-service/    # Payment service
├── internal/               # Private application code
[v]│   ├── auth/               # Authentication logic
│   ├── cache/              # Redis cache management
│   ├── config/             # Configuration management
│   ├── database/           # Database connection
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # HTTP middleware
│   ├── models/             # Data models and structures
│   ├── queue/              # Message queue management
│   ├── repository/         # Data access layer
│   └── service/            # Business logic layer
├── pkg/                    # Public/shared libraries
│   ├── api/                # API definitions and contracts
│   ├── errors/             # Error definitions and handling
[v]│   ├── jwt/                # JWT token management
│   ├── logger/             # Logging utilities
│   ├── pagination/         # Pagination helpers
│   ├── response/           # Standard response formats
│   ├── utils/              # Utility functions
│   └── validation/         # Input validation utilities
├── services/               # Microservice implementations
│   ├── booking/            # Booking service implementation
│   │   ├── handlers/       # HTTP handlers
│   │   ├── repository/     # Data access layer
│   │   └── service/        # Business logic
│   ├── flight/             # Flight service implementation
│   │   ├── handlers/       # HTTP handlers
│   │   ├── repository/     # Data access layer
│   │   └── service/        # Business logic
│   ├── notification/       # Notification service implementation
│   │   ├── handlers/       # HTTP handlers
│   │   ├── repository/     # Data access layer
│   │   └── service/        # Business logic
│   └── payment/            # Payment service implementation
│       ├── handlers/       # HTTP handlers
│       ├── repository/     # Data access layer
│       └── service/        # Business logic
├── migrations/             # PostgreSQL migration files
├── seeds/                  # Demo data and seeding scripts
├── deployments/            # Deployment configurations
├── docs/                   # Documentation
├── tests/                  # Test files
├── web/                    # Web assets
│   ├── static/             # Static files
│   └── templates/          # HTML templates
├── scripts/                # Build and automation scripts
├── tools/                  # Development tools
├── configs/                # Configuration files
├── go.mod                  # Go module definition
├── docker-compose.yml      # Docker services configuration
├── Dockerfile              # Container build configuration
├── Makefile                # Build and automation commands
└── db.dbml                 # Database schema definition
```

## Detailed Directory Structure

### `/cmd` - Application Entry Points

Contains the main applications for this project. Each subdirectory represents a service entry point.

- `api-gateway/` - API Gateway service main entry point
- `flight-service/` - Flight service main entry point
- `booking-service/` - Booking service main entry point
- `payment-service/` - Payment service main entry point
- `notification-service/` - Notification service main entry point

### `/internal` - Private Application Code

Contains private application code that shouldn't be imported by other applications.

- `models/` - Data models and structures
- `config/` - Configuration management
- `database/` - Database connection and management
- `cache/` - Redis cache management
- `queue/` - Message queue management
- `auth/` - Authentication logic
- `middleware/` - HTTP middleware
- `handlers/` - HTTP request handlers
- `repository/` - Data access layer
- `service/` - Business logic layer

### `/pkg` - Public/Shared Libraries

Contains code that's safe to use by external applications.

- `api/` - API definitions and contracts
- `utils/` - Utility functions
- `errors/` - Error definitions and handling
- `logger/` - Logging utilities
- `jwt/` - JWT token management
- `validation/` - Input validation utilities
- `pagination/` - Pagination helpers
- `response/` - Standard response formats

### `/services` - Microservice Implementations

Each subdirectory contains a complete microservice implementation.

#### `/services/flight` - Flight Management Service

- `handlers/` - HTTP handlers for flight operations
- `repository/` - Flight data access layer
- `service/` - Flight business logic

#### `/services/booking` - Booking Management Service

- `handlers/` - HTTP handlers for booking operations
- `repository/` - Booking data access layer
- `service/` - Booking business logic (including overbooking)

#### `/services/payment` - Payment Processing Service

- `handlers/` - HTTP handlers for payment operations
- `repository/` - Payment data access layer
- `service/` - Payment processing logic

#### `/services/notification` - Notification Service

- `handlers/` - HTTP handlers for notification operations
- `repository/` - Notification data access layer
- `service/` - Email/SMS notification logic

### `/migrations` - Database Migrations

Contains PostgreSQL migration files for database schema management.

- `001_create_core_entities.up.sql` - Create core entities (airlines, airports, aircraft)
- `002_create_flights.up.sql` - Create flights and schedules tables
- `003_create_inventory.up.sql` - Create flight inventory table
- `004_create_users_bookings.up.sql` - Create users and booking tables
- `005_create_overbooking.up.sql` - Create overbooking management tables
- `006_create_analytics.up.sql` - Create analytics and search optimization tables
- `007_create_seats_system.up.sql` - Create seat management and system tables
- `008_insert_seed_data.up.sql` - Insert demo/seed data
- `*.down.sql` - Corresponding rollback migrations

### `/seeds` - Demo Data

Contains SQL files and scripts for populating the database with demo data.

- `01_users.sql` through `09_system_config.sql` - Structured demo data
- `run_seeds.sh` - Intelligent seeding script with validation
- `README.md` - Seeding documentation and usage guide

### `/deployments` - Deployment Configurations

Contains deployment and orchestration configurations.

- Kubernetes manifests
- Docker Compose files
- Prometheus configuration
- Environment-specific configurations

### `/docs` - Documentation

Contains project documentation.

- API documentation
- Architecture diagrams
- Development guides
- Deployment instructions

### `/tests` - Test Files

Contains test files and test data.

- Unit tests
- Integration tests
- Performance tests
- Test utilities

### `/web` - Web Assets

Contains web application assets.

- `static/` - Static files (CSS, JS, images)
- `templates/` - HTML templates

### `/scripts` - Build and Automation Scripts

Contains build, deployment, and automation scripts.

- Build scripts
- CI/CD scripts
- Utility scripts

### `/tools` - Development Tools

Contains tools needed for development.

- Code generation tools
- Database tools
- Testing tools

### `/configs` - Configuration Files

Contains configuration files for different environments.

- Development configuration
- Production configuration
- Test configuration

## Additional Files in Project Root

### Core Configuration Files

- **`go.mod`** - Go module definition and dependencies
- **`docker-compose.yml`** - Multi-service Docker orchestration with PostgreSQL, Redis, RabbitMQ, Elasticsearch, monitoring stack
- **`Dockerfile`** - Container build configuration for services
- **`Makefile`** - Build automation, database operations, and development commands

### Documentation Files

- **`db.dbml`** - Database schema definition in DBML format for PostgreSQL
- **`folder_structure.md`** - This file, documenting the project structure
- **`target.md`** - Project requirements and specifications
- **`flight_booking_todo.md`** - Detailed project development roadmap
- **`booking_sequence.mermaid`** - Booking process sequence diagram
- **`flight_query_sequence.mermaid`** - Flight search sequence diagram

### Development Files

- **`.gitignore`** - Git ignore patterns for Go projects

## Architecture Notes

This structure follows the **Clean Architecture** principles:

1. **Separation of Concerns**: Each layer has a specific responsibility
2. **Dependency Inversion**: Inner layers don't depend on outer layers
3. **Interface Segregation**: Well-defined interfaces between layers
4. **Single Responsibility**: Each package has a single, well-defined purpose

## Microservices Architecture

The system is designed as a microservices architecture with:

- **API Gateway**: Single entry point, handles routing, authentication, rate limiting
- **Flight Service**: Manages flight data, schedules, and inventory
- **Booking Service**: Handles booking lifecycle, overbooking logic
- **Payment Service**: Processes payments and handles payment states
- **Notification Service**: Sends emails, SMS, and push notifications

## High Concurrency Design

The structure supports high concurrency through:

- **Redis Caching**: Fast data access for hot data
- **Message Queues**: Asynchronous processing
- **Database Sharding**: Horizontal scaling capability
- **Read/Write Separation**: Optimized database access patterns
- **Distributed Locking**: Prevents race conditions in critical sections
