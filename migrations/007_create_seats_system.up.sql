-- Create seat_maps table
CREATE TABLE seat_maps (
    id BIGSERIAL PRIMARY KEY,
    aircraft_id INTEGER NOT NULL REFERENCES aircraft(id),
    seat_number VARCHAR(10) NOT NULL,
    class_type class_type NOT NULL,
    seat_type seat_type NOT NULL,
    row_number INTEGER NOT NULL,
    column_letter VARCHAR(1) NOT NULL,
    is_exit_row BOOLEAN DEFAULT FALSE,
    extra_legroom BOOLEAN DEFAULT FALSE,
    blocked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint for aircraft and seat combination
    CONSTRAINT uk_aircraft_seat UNIQUE (aircraft_id, seat_number)
);

-- Create indexes for seat_maps
CREATE INDEX idx_seat_maps_aircraft_id ON seat_maps(aircraft_id);
CREATE INDEX idx_seat_maps_class_type ON seat_maps(class_type);
CREATE INDEX idx_seat_maps_seat_type ON seat_maps(seat_type);

-- Add comments
COMMENT ON TABLE seat_maps IS 'Aircraft seat layout definitions';
COMMENT ON COLUMN seat_maps.seat_number IS 'Seat identifier (e.g., 12A)';
COMMENT ON COLUMN seat_maps.blocked IS 'Permanently blocked seat';

-- Create seat_assignments table
CREATE TABLE seat_assignments (
    id BIGSERIAL PRIMARY KEY,
    schedule_id BIGINT NOT NULL REFERENCES flight_schedules(id),
    seat_id BIGINT NOT NULL REFERENCES seat_maps(id),
    booking_id BIGINT REFERENCES bookings(id),
    passenger_id BIGINT REFERENCES booking_passengers(id),
    status seat_status DEFAULT 'AVAILABLE',
    reserved_at TIMESTAMP,
    confirmed_at TIMESTAMP,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint for schedule and seat combination
    CONSTRAINT uk_schedule_seat UNIQUE (schedule_id, seat_id)
);

-- Create trigger for seat_assignments updated_at
CREATE TRIGGER update_seat_assignments_updated_at 
    BEFORE UPDATE ON seat_assignments 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for seat_assignments
CREATE INDEX idx_seat_assignments_schedule_id ON seat_assignments(schedule_id);
CREATE INDEX idx_seat_assignments_seat_id ON seat_assignments(seat_id);
CREATE INDEX idx_seat_assignments_booking_id ON seat_assignments(booking_id);
CREATE INDEX idx_seat_assignments_passenger_id ON seat_assignments(passenger_id);
CREATE INDEX idx_seat_assignments_status ON seat_assignments(status);
CREATE INDEX idx_seat_assignments_schedule_status ON seat_assignments(schedule_id, status);

-- Add comments
COMMENT ON TABLE seat_assignments IS 'Seat assignments for flight schedules';
COMMENT ON COLUMN seat_assignments.booking_id IS 'NULL if seat is blocked/maintenance';
COMMENT ON COLUMN seat_assignments.reserved_at IS 'When seat was reserved';
COMMENT ON COLUMN seat_assignments.confirmed_at IS 'When seat was confirmed';
COMMENT ON COLUMN seat_assignments.version IS 'Optimistic locking version';

-- Create system_settings table
CREATE TABLE system_settings (
    id SERIAL PRIMARY KEY,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT NOT NULL,
    setting_type setting_type DEFAULT 'STRING',
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for system_settings updated_at
CREATE TRIGGER update_system_settings_updated_at 
    BEFORE UPDATE ON system_settings 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for system_settings
CREATE INDEX idx_system_settings_setting_key ON system_settings(setting_key);
CREATE INDEX idx_system_settings_is_public ON system_settings(is_public);

-- Add comments
COMMENT ON TABLE system_settings IS 'Application configuration settings';
COMMENT ON COLUMN system_settings.is_public IS 'Whether setting can be exposed to frontend';

-- Create api_rate_limits table (TTL table)
CREATE TABLE api_rate_limits (
    id BIGSERIAL PRIMARY KEY,
    identifier VARCHAR(100) NOT NULL,
    endpoint VARCHAR(200) NOT NULL,
    request_count INTEGER DEFAULT 1,
    window_start TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- Create indexes for api_rate_limits
CREATE INDEX idx_api_rate_limits_rate_limit ON api_rate_limits(identifier, endpoint);
CREATE INDEX idx_api_rate_limits_expires_at ON api_rate_limits(expires_at);

-- Add comments
COMMENT ON TABLE api_rate_limits IS 'API rate limiting records';
COMMENT ON COLUMN api_rate_limits.identifier IS 'User ID, IP address, or API key';