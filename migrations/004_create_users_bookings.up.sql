-- Create users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    active BOOLEAN DEFAULT TRUE,
    tier VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for users updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for users
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_active ON users(active);

-- Add comments
COMMENT ON TABLE users IS 'System users for JWT authentication';

-- Create bookings table
CREATE TABLE bookings (
    id BIGSERIAL PRIMARY KEY,
    booking_uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    schedule_id BIGINT NOT NULL REFERENCES flight_schedules(id),
    class_type class_type NOT NULL,
    passenger_count SMALLINT NOT NULL DEFAULT 1,
    total_amount DECIMAL(12,2) NOT NULL,
    status booking_status DEFAULT 'PENDING',
    seat_numbers JSONB,
    special_requests TEXT,
    booking_source VARCHAR(20) DEFAULT 'WEB',
    expires_at TIMESTAMP,
    confirmed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for bookings updated_at
CREATE TRIGGER update_bookings_updated_at 
    BEFORE UPDATE ON bookings 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for bookings
CREATE INDEX idx_bookings_booking_uuid ON bookings(booking_uuid);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_schedule_id ON bookings(schedule_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_expires_at ON bookings(expires_at);
CREATE INDEX idx_bookings_user_status ON bookings(user_id, status);
CREATE INDEX idx_bookings_schedule_status ON bookings(schedule_id, status);
CREATE INDEX idx_bookings_created_at ON bookings(created_at);

-- Add comments
COMMENT ON TABLE bookings IS 'Flight booking records';
COMMENT ON COLUMN bookings.booking_uuid IS 'Public booking identifier';
COMMENT ON COLUMN bookings.passenger_count IS 'Number of passengers';
COMMENT ON COLUMN bookings.total_amount IS 'Total booking amount';
COMMENT ON COLUMN bookings.seat_numbers IS 'Assigned seat numbers';
COMMENT ON COLUMN bookings.special_requests IS 'Special passenger requests';
COMMENT ON COLUMN bookings.booking_source IS 'Booking channel';
COMMENT ON COLUMN bookings.expires_at IS 'Reservation expiry time';
COMMENT ON COLUMN bookings.confirmed_at IS 'Booking confirmation time';
COMMENT ON COLUMN bookings.cancelled_at IS 'Cancellation time';

-- Create booking_passengers table
CREATE TABLE booking_passengers (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    date_of_birth DATE NOT NULL,
    passport_number VARCHAR(20),
    nationality VARCHAR(2),
    seat_number VARCHAR(10),
    meal_preference VARCHAR(20),
    special_needs TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for booking_passengers
CREATE INDEX idx_booking_passengers_booking_id ON booking_passengers(booking_id);
CREATE INDEX idx_booking_passengers_passport ON booking_passengers(passport_number);

-- Add comments
COMMENT ON TABLE booking_passengers IS 'Passenger details for bookings';
COMMENT ON COLUMN booking_passengers.nationality IS 'Country code';
COMMENT ON COLUMN booking_passengers.seat_number IS 'Assigned seat';

-- Create booking_payments table
CREATE TABLE booking_payments (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id),
    payment_method payment_method NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    payment_status payment_status DEFAULT 'PENDING',
    transaction_id VARCHAR(100) UNIQUE,
    gateway_response JSONB,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for booking_payments
CREATE INDEX idx_booking_payments_booking_id ON booking_payments(booking_id);
CREATE INDEX idx_booking_payments_status ON booking_payments(payment_status);
CREATE INDEX idx_booking_payments_transaction_id ON booking_payments(transaction_id);
CREATE INDEX idx_booking_payments_processed_at ON booking_payments(processed_at);

-- Add comments
COMMENT ON TABLE booking_payments IS 'Payment records for bookings';
COMMENT ON COLUMN booking_payments.transaction_id IS 'External payment gateway transaction ID';
COMMENT ON COLUMN booking_payments.gateway_response IS 'Payment gateway response';
COMMENT ON COLUMN booking_payments.processed_at IS 'Payment processing time';

-- Create booking_status_log table
CREATE TABLE booking_status_log (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id),
    old_status VARCHAR(20),
    new_status VARCHAR(20) NOT NULL,
    reason TEXT,
    changed_by VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for booking_status_log
CREATE INDEX idx_booking_status_log_booking_id ON booking_status_log(booking_id);
CREATE INDEX idx_booking_status_log_new_status ON booking_status_log(new_status);
CREATE INDEX idx_booking_status_log_created_at ON booking_status_log(created_at);

-- Add comments
COMMENT ON TABLE booking_status_log IS 'Audit log for booking status changes';
COMMENT ON COLUMN booking_status_log.changed_by IS 'System or user identifier';