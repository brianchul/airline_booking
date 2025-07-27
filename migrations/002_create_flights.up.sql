-- Create flights table
CREATE TABLE flights (
    id SERIAL PRIMARY KEY,
    flight_number VARCHAR(10) NOT NULL,
    airline_id INTEGER NOT NULL REFERENCES airlines(id),
    departure_airport_id INTEGER NOT NULL REFERENCES airports(id),
    arrival_airport_id INTEGER NOT NULL REFERENCES airports(id),
    aircraft_id INTEGER NOT NULL REFERENCES aircraft(id),
    base_price_economy DECIMAL(10,2) NOT NULL,
    base_price_business DECIMAL(10,2),
    base_price_first DECIMAL(10,2),
    duration_minutes INTEGER NOT NULL,
    distance_km INTEGER,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for flights updated_at
CREATE TRIGGER update_flights_updated_at 
    BEFORE UPDATE ON flights 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for flights
CREATE INDEX idx_flights_flight_number ON flights(flight_number);
CREATE INDEX idx_flights_airline_id ON flights(airline_id);
CREATE INDEX idx_flights_route ON flights(departure_airport_id, arrival_airport_id);
CREATE INDEX idx_flights_departure_airport ON flights(departure_airport_id);
CREATE INDEX idx_flights_arrival_airport ON flights(arrival_airport_id);
CREATE INDEX idx_flights_aircraft_id ON flights(aircraft_id);
CREATE INDEX idx_flights_active ON flights(active);
CREATE UNIQUE INDEX uk_flight_airline ON flights(flight_number, airline_id);

-- Add comments
COMMENT ON TABLE flights IS 'Flight route definitions';
COMMENT ON COLUMN flights.flight_number IS 'Flight number';
COMMENT ON COLUMN flights.base_price_economy IS 'Base price for economy class';
COMMENT ON COLUMN flights.base_price_business IS 'Base price for business class';
COMMENT ON COLUMN flights.base_price_first IS 'Base price for first class';
COMMENT ON COLUMN flights.duration_minutes IS 'Flight duration in minutes';
COMMENT ON COLUMN flights.distance_km IS 'Flight distance in kilometers';

-- Create flight_schedules table with partitioning
CREATE TABLE flight_schedules (
    id BIGSERIAL PRIMARY KEY,
    flight_id INTEGER NOT NULL REFERENCES flights(id),
    departure_time TIMESTAMP NOT NULL,
    arrival_time TIMESTAMP NOT NULL,
    status flight_status DEFAULT 'SCHEDULED',
    gate VARCHAR(10),
    terminal VARCHAR(10),
    actual_departure_time TIMESTAMP,
    actual_arrival_time TIMESTAMP,
    delay_minutes INTEGER DEFAULT 0,
    cancellation_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (departure_time);

-- Create trigger for flight_schedules updated_at
CREATE TRIGGER update_flight_schedules_updated_at 
    BEFORE UPDATE ON flight_schedules 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create partitions for flight_schedules (current year and next year)
CREATE TABLE flight_schedules_2024 PARTITION OF flight_schedules
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE flight_schedules_2025 PARTITION OF flight_schedules
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

-- Create indexes for flight_schedules
CREATE INDEX idx_flight_schedules_flight_id ON flight_schedules(flight_id);
CREATE INDEX idx_flight_schedules_departure_time ON flight_schedules(departure_time);
CREATE INDEX idx_flight_schedules_arrival_time ON flight_schedules(arrival_time);
CREATE INDEX idx_flight_schedules_status ON flight_schedules(status);
CREATE INDEX idx_flight_schedules_date_range ON flight_schedules(departure_time, arrival_time);
CREATE INDEX idx_flight_schedules_search ON flight_schedules(flight_id, departure_time, status);

-- Add comments
COMMENT ON TABLE flight_schedules IS 'Scheduled flight instances';
COMMENT ON COLUMN flight_schedules.departure_time IS 'Scheduled departure time';
COMMENT ON COLUMN flight_schedules.arrival_time IS 'Scheduled arrival time';
COMMENT ON COLUMN flight_schedules.gate IS 'Departure gate';
COMMENT ON COLUMN flight_schedules.terminal IS 'Departure terminal';
COMMENT ON COLUMN flight_schedules.actual_departure_time IS 'Actual departure time';
COMMENT ON COLUMN flight_schedules.actual_arrival_time IS 'Actual arrival time';
COMMENT ON COLUMN flight_schedules.delay_minutes IS 'Delay in minutes';
COMMENT ON COLUMN flight_schedules.cancellation_reason IS 'Reason for cancellation';