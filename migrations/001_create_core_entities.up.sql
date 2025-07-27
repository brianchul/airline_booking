-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom types for enums
CREATE TYPE flight_status AS ENUM ('SCHEDULED', 'DELAYED', 'CANCELLED', 'DEPARTED', 'ARRIVED');
CREATE TYPE class_type AS ENUM ('ECONOMY', 'BUSINESS', 'FIRST');
CREATE TYPE booking_status AS ENUM ('PENDING', 'QUEUED', 'PROCESSING', 'RESERVED', 'CONFIRMED', 'CANCELLED', 'EXPIRED', 'REFUNDED');
CREATE TYPE payment_method AS ENUM ('CREDIT_CARD', 'DEBIT_CARD', 'BANK_TRANSFER', 'WALLET');
CREATE TYPE payment_status AS ENUM ('PENDING', 'PROCESSING', 'COMPLETED', 'FAILED', 'REFUNDED');
CREATE TYPE route_type AS ENUM ('DOMESTIC', 'INTERNATIONAL', 'ALL');
CREATE TYPE overbooking_class_type AS ENUM ('ECONOMY', 'BUSINESS', 'FIRST', 'ALL');
CREATE TYPE compensation_type AS ENUM ('CASH', 'VOUCHER', 'MILES', 'UPGRADE');
CREATE TYPE flight_type AS ENUM ('DOMESTIC', 'INTERNATIONAL');
CREATE TYPE seat_type AS ENUM ('WINDOW', 'AISLE', 'MIDDLE');
CREATE TYPE seat_status AS ENUM ('AVAILABLE', 'RESERVED', 'CONFIRMED', 'BLOCKED');
CREATE TYPE resolution_status AS ENUM ('RESOLVED', 'PENDING', 'ESCALATED');
CREATE TYPE setting_type AS ENUM ('STRING', 'INTEGER', 'BOOLEAN', 'JSON');

-- Create airlines table
CREATE TABLE airlines (
    id SERIAL PRIMARY KEY,
    code VARCHAR(3) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    country VARCHAR(2) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create function for updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for airlines updated_at
CREATE TRIGGER update_airlines_updated_at 
    BEFORE UPDATE ON airlines 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for airlines
CREATE INDEX idx_airlines_code ON airlines(code);
CREATE INDEX idx_airlines_country ON airlines(country);
CREATE INDEX idx_airlines_active ON airlines(active);

-- Add comments
COMMENT ON TABLE airlines IS 'Airline companies';
COMMENT ON COLUMN airlines.code IS 'IATA airline code';
COMMENT ON COLUMN airlines.name IS 'Airline name';
COMMENT ON COLUMN airlines.country IS 'Country code';

-- Create airports table
CREATE TABLE airports (
    id SERIAL PRIMARY KEY,
    code VARCHAR(3) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    city VARCHAR(50) NOT NULL,
    country VARCHAR(2) NOT NULL,
    timezone VARCHAR(50) NOT NULL,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for airports updated_at
CREATE TRIGGER update_airports_updated_at 
    BEFORE UPDATE ON airports 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for airports
CREATE INDEX idx_airports_code ON airports(code);
CREATE INDEX idx_airports_city ON airports(city);
CREATE INDEX idx_airports_country ON airports(country);
CREATE INDEX idx_airports_active ON airports(active);

-- Add comments
COMMENT ON TABLE airports IS 'Airport information';
COMMENT ON COLUMN airports.code IS 'IATA airport code';
COMMENT ON COLUMN airports.name IS 'Airport name';
COMMENT ON COLUMN airports.city IS 'City name';
COMMENT ON COLUMN airports.country IS 'Country code';
COMMENT ON COLUMN airports.timezone IS 'Timezone identifier';
COMMENT ON COLUMN airports.latitude IS 'Latitude coordinate';
COMMENT ON COLUMN airports.longitude IS 'Longitude coordinate';

-- Create aircraft table
CREATE TABLE aircraft (
    id SERIAL PRIMARY KEY,
    model VARCHAR(50) NOT NULL,
    manufacturer VARCHAR(50) NOT NULL,
    capacity_economy INTEGER NOT NULL,
    capacity_business INTEGER DEFAULT 0,
    capacity_first INTEGER DEFAULT 0,
    total_capacity INTEGER GENERATED ALWAYS AS (capacity_economy + capacity_business + capacity_first) STORED,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for aircraft updated_at
CREATE TRIGGER update_aircraft_updated_at 
    BEFORE UPDATE ON aircraft 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for aircraft
CREATE INDEX idx_aircraft_model ON aircraft(model);
CREATE INDEX idx_aircraft_manufacturer ON aircraft(manufacturer);
CREATE INDEX idx_aircraft_total_capacity ON aircraft(total_capacity);
CREATE INDEX idx_aircraft_active ON aircraft(active);

-- Add comments
COMMENT ON TABLE aircraft IS 'Aircraft models and configurations';
COMMENT ON COLUMN aircraft.model IS 'Aircraft model';
COMMENT ON COLUMN aircraft.manufacturer IS 'Aircraft manufacturer';
COMMENT ON COLUMN aircraft.capacity_economy IS 'Economy class capacity';
COMMENT ON COLUMN aircraft.capacity_business IS 'Business class capacity';
COMMENT ON COLUMN aircraft.capacity_first IS 'First class capacity';