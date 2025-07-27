# Database Seeding for Airline Booking System

This directory contains scripts and data files to populate the database with realistic demo data for testing and development.

## Overview

The seeding process creates a comprehensive dataset that includes:

- **6 Demo Users** with JWT authentication (password: `password123`)
- **21 Airlines** from major regions worldwide
- **31 Airports** covering major international hubs
- **23 Aircraft Types** with realistic configurations
- **24 Flight Routes** between popular destinations
- **Flight Schedules** for the next 30 days
- **Flight Inventory** with realistic seat availability
- **Sample Bookings** in various states
- **System Configuration** and business rules

## Quick Start

### Using Makefile (Recommended)

```bash
# Setup complete demo environment
make demo-setup

# Or seed existing database
make seed-db

# Reset and reseed database
make seed-reset

# Check database status
make db-status
```

### Manual Execution

```bash
# Run the seeding script directly
./seeds/run_seeds.sh

# With custom database parameters
./seeds/run_seeds.sh --host db.example.com --port 5433
```

## Seed Files

### Core Data
- `01_users.sql` - Demo user accounts for JWT authentication
- `02_airlines.sql` - Major airlines worldwide
- `03_airports.sql` - International airports with coordinates
- `04_aircraft.sql` - Aircraft types and seat configurations

### Flight Data
- `05_flights.sql` - Popular flight routes with realistic pricing
- `06_flight_schedules.sql` - Flight schedules for next 30 days
- `07_flight_inventory.sql` - Seat availability and dynamic pricing

### Demo Scenarios
- `08_sample_bookings.sql` - Sample bookings in various states
- `09_system_config.sql` - Business rules and system settings

## Demo User Accounts

All users have the password: `password123`

| Username | Email | Role | Description |
|----------|-------|------|-------------|
| `admin` | admin@airline.com | Admin | System administrator |
| `john.doe` | john.doe@example.com | User | Regular customer |
| `jane.smith` | jane.smith@example.com | User | Regular customer |
| `mike.wilson` | mike.wilson@example.com | User | Regular customer |
| `sarah.johnson` | sarah.johnson@example.com | User | Regular customer |
| `demo.user` | demo@example.com | User | Demo/testing account |

## Data Highlights

### Flight Routes
Popular international routes with realistic pricing:
- **JFK ↔ LAX** - Domestic US transcontinental
- **LHR ↔ JFK** - Transatlantic premium route
- **SIN ↔ LHR** - Long-haul Asia-Europe
- **DXB ↔ JFK** - Middle East hub connections
- **NRT ↔ LAX** - Transpacific route

### Pricing Strategy
- **Economy**: $180 - $1,400 depending on route and demand
- **Business**: $520 - $4,000 with premium service
- **First**: $2,500 - $8,000 for luxury travel
- **Dynamic pricing** based on availability and time to departure

### Booking Scenarios
Sample bookings demonstrate various states:
- **Confirmed** bookings with payments
- **Reserved** bookings awaiting payment
- **Expired** bookings for timeout scenarios
- **Cancelled** bookings for refund testing
- **Pending** bookings in queue

### Overbooking Strategy
Realistic overbooking rules by airline and route type:
- **US Domestic**: 6-12% overbooking rates
- **International**: 3-8% more conservative rates
- **Business/First**: Lower overbooking (2-5%)
- **Regional variations** for EU, US, Asia regulations

## System Configuration

### JWT Settings
- **Secret Key**: `your-super-secret-jwt-key-change-in-production`
- **Token Expiry**: 24 hours
- **Refresh Token**: 7 days

### Rate Limiting
- **Search API**: 20 requests/minute
- **Booking API**: 5 requests/minute
- **Auth API**: 10 requests/minute

### Cache Settings
- **Search Results**: 5 minutes TTL
- **Inventory**: 2 minutes TTL
- **Flight Details**: 10 minutes TTL

### Business Rules
- **Booking Expiry**: 15 minutes for payment
- **Seat Hold**: 10 minutes during selection
- **Advance Booking**: 2 hours minimum, 365 days maximum

## Customization

### Environment Variables
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=airline_booking
export DB_USER=postgres
export DB_PASSWORD=rootpass
```

### Script Options
```bash
./run_seeds.sh --help
```

Available options:
- `--host HOST` - Database host
- `--port PORT` - Database port  
- `--database DATABASE` - Database name
- `--user USER` - Database user
- `--password PASSWORD` - Database password

## Verification

After seeding, verify the data:

```sql
-- Check record counts
SELECT 'Users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Airlines', COUNT(*) FROM airlines
UNION ALL  
SELECT 'Airports', COUNT(*) FROM airports
UNION ALL
SELECT 'Flights', COUNT(*) FROM flights
UNION ALL
SELECT 'Schedules', COUNT(*) FROM flight_schedules
UNION ALL
SELECT 'Inventory', COUNT(*) FROM flight_inventory
UNION ALL
SELECT 'Bookings', COUNT(*) FROM bookings;

-- Check upcoming flights
SELECT f.flight_number, al.name as airline, 
       dep.code as from_airport, arr.code as to_airport,
       fs.departure_time, fs.status
FROM flight_schedules fs
JOIN flights f ON fs.flight_id = f.id
JOIN airlines al ON f.airline_id = al.id
JOIN airports dep ON f.departure_airport_id = dep.id
JOIN airports arr ON f.arrival_airport_id = arr.id
WHERE fs.departure_time > CURRENT_TIMESTAMP
ORDER BY fs.departure_time
LIMIT 10;
```

## Troubleshooting

### Common Issues

1. **Connection Failed**
   ```bash
   # Check if PostgreSQL is running
   docker-compose ps postgres-master
   
   # Check connection
   psql -h localhost -p 5432 -U postgres -d airline_booking -c "SELECT 1;"
   ```

2. **Tables Not Found**
   ```bash
   # Run migrations first
   make migrate-up
   ```

3. **Permission Denied**
   ```bash
   # Make script executable
   chmod +x seeds/run_seeds.sh
   ```

4. **Duplicate Data**
   - The seed script uses `ON CONFLICT DO NOTHING` to prevent duplicates
   - Use `make seed-reset` to start fresh

### Reset Everything
```bash
# Complete reset
make dev-down
make dev-up
sleep 10
make migrate-up
make seed-db
```

## Production Notes

⚠️ **Important**: This seed data is for development/demo only:

- Default passwords are weak (`password123`)
- JWT secret key should be changed
- Realistic but not production-scale data volumes
- Some business rules are simplified for demo purposes

For production deployment:
- Generate secure passwords and JWT secrets
- Implement proper user registration
- Scale data volumes appropriately  
- Review and adjust business rules
- Implement proper backup/restore procedures