#!/bin/bash

# Airline Booking System - Database Seeding Script
# This script populates the database with demo data for testing

set -e

# Docker container parameters
CONTAINER_NAME=${CONTAINER_NAME:-airline_booking-postgres-master-1}
DB_NAME=${DB_NAME:-airline_booking}
DB_USER=${DB_USER:-postgres}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to execute SQL file
execute_sql_file() {
    local file=$1
    local description=$2

    print_status "Executing $description..."

    if [ ! -f "$file" ]; then
        print_error "File $file not found!"
        exit 1
    fi

    if docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME < "$file" > /dev/null 2>&1; then
        print_success "$description completed"
    else
        print_error "Failed to execute $description"
        exit 1
    fi
}

# Function to check database connection
check_database_connection() {
    print_status "Checking database connection..."

    if docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
        print_success "Database connection successful"
    else
        print_error "Cannot connect to database. Please check your Docker container and parameters."
        print_error "Container: $CONTAINER_NAME, Database: $DB_NAME, User: $DB_USER"
        exit 1
    fi
}

# Function to check if tables exist
check_tables_exist() {
    print_status "Checking if database schema exists..."

    local table_count=$(docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name IN ('airlines', 'airports', 'flights');" 2>/dev/null || echo "0")

    if [ "$table_count" -lt 3 ]; then
        print_error "Database schema not found. Please run migrations first:"
        print_error "make migrate-up"
        exit 1
    else
        print_success "Database schema found"
    fi
}

# Function to check if data already exists
check_existing_data() {
    print_status "Checking for existing data..."

    local user_count=$(docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM users;" 2>/dev/null || echo "0")

    if [ "$user_count" -gt 0 ]; then
        print_warning "Data already exists in the database."
        read -p "Do you want to continue and potentially add duplicate data? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Seeding cancelled by user"
            exit 0
        fi
    fi
}

# Main seeding function
seed_database() {
    local seed_dir="$(dirname "$0")"

    print_status "Starting database seeding process..."
    print_status "Seed directory: $seed_dir"

    # Execute seed files in order
    execute_sql_file "$seed_dir/01_users.sql" "User accounts"
    execute_sql_file "$seed_dir/02_airlines.sql" "Airlines data"
    execute_sql_file "$seed_dir/03_airports.sql" "Airports data"
    execute_sql_file "$seed_dir/04_aircraft.sql" "Aircraft data"
    execute_sql_file "$seed_dir/05_flights.sql" "Flight routes"
    execute_sql_file "$seed_dir/06_flight_schedules.sql" "Flight schedules"
    execute_sql_file "$seed_dir/07_flight_inventory.sql" "Flight inventory"
    execute_sql_file "$seed_dir/08_sample_bookings.sql" "Sample bookings"
    execute_sql_file "$seed_dir/09_system_config.sql" "System configuration"

    print_success "Database seeding completed successfully!"
}

# Function to display summary
display_summary() {
    print_status "Seeding Summary:"

    local counts=$(docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -t -c "
        SELECT
            'Users: ' || (SELECT COUNT(*) FROM users) || E'\n' ||
            'Airlines: ' || (SELECT COUNT(*) FROM airlines) || E'\n' ||
            'Airports: ' || (SELECT COUNT(*) FROM airports) || E'\n' ||
            'Aircraft: ' || (SELECT COUNT(*) FROM aircraft) || E'\n' ||
            'Flights: ' || (SELECT COUNT(*) FROM flights) || E'\n' ||
            'Flight Schedules: ' || (SELECT COUNT(*) FROM flight_schedules) || E'\n' ||
            'Flight Inventory: ' || (SELECT COUNT(*) FROM flight_inventory) || E'\n' ||
            'Bookings: ' || (SELECT COUNT(*) FROM bookings) || E'\n' ||
            'System Settings: ' || (SELECT COUNT(*) FROM system_settings);
    " 2>/dev/null)

    echo "$counts"

    print_success "Demo data is ready for testing!"
    echo
    print_status "Demo user accounts (password: password123):"
    echo "  - admin / admin@airline.com"
    echo "  - john.doe / john.doe@example.com"
    echo "  - jane.smith / jane.smith@example.com"
    echo "  - demo.user / demo@example.com"
}

# Function to display help
show_help() {
    echo "Airline Booking System - Database Seeding Script"
    echo
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  --container CONTAINER   Docker container name (default: airline_booking-postgres-master-1)"
    echo "  --database DATABASE     Database name (default: airline_booking)"
    echo "  --user USER             Database user (default: postgres)"
    echo
    echo "Environment Variables:"
    echo "  CONTAINER_NAME          Docker container name"
    echo "  DB_NAME                 Database name"
    echo "  DB_USER                 Database user"
    echo
    echo "Examples:"
    echo "  $0                                    # Use default container"
    echo "  $0 --container my-db-container        # Custom container name"
    echo "  CONTAINER_NAME=my-db $0               # Custom container via env var"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        --container)
            CONTAINER_NAME="$2"
            shift 2
            ;;
        --database)
            DB_NAME="$2"
            shift 2
            ;;
        --user)
            DB_USER="$2"
            shift 2
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Main execution
main() {
    echo "=============================================="
    echo "  Airline Booking System - Database Seeding"
    echo "=============================================="
    echo

    check_database_connection
    check_tables_exist
    check_existing_data
    seed_database
    display_summary

    echo
    print_success "Seeding process completed successfully!"
}

# Run main function
main