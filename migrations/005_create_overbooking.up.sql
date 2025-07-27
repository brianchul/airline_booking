-- Create overbooking_rules table
CREATE TABLE overbooking_rules (
    id SERIAL PRIMARY KEY,
    airline_id INTEGER NOT NULL REFERENCES airlines(id),
    aircraft_id INTEGER REFERENCES aircraft(id),
    route_type route_type DEFAULT 'ALL',
    class_type overbooking_class_type DEFAULT 'ALL',
    base_overbooking_rate DECIMAL(5,2) NOT NULL,
    max_overbooking_rate DECIMAL(5,2) NOT NULL,
    no_show_rate DECIMAL(5,2) NOT NULL,
    seasonal_factor DECIMAL(3,2) DEFAULT 1.00,
    active BOOLEAN DEFAULT TRUE,
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for overbooking_rules updated_at
CREATE TRIGGER update_overbooking_rules_updated_at 
    BEFORE UPDATE ON overbooking_rules 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for overbooking_rules
CREATE INDEX idx_overbooking_rules_airline_id ON overbooking_rules(airline_id);
CREATE INDEX idx_overbooking_rules_aircraft_id ON overbooking_rules(aircraft_id);
CREATE INDEX idx_overbooking_rules_effective ON overbooking_rules(effective_from, effective_to);
CREATE INDEX idx_overbooking_rules_active ON overbooking_rules(active);

-- Add comments
COMMENT ON TABLE overbooking_rules IS 'Overbooking strategy rules';
COMMENT ON COLUMN overbooking_rules.aircraft_id IS 'Specific aircraft or NULL for all';
COMMENT ON COLUMN overbooking_rules.base_overbooking_rate IS 'Base overbooking percentage';
COMMENT ON COLUMN overbooking_rules.max_overbooking_rate IS 'Maximum overbooking percentage';
COMMENT ON COLUMN overbooking_rules.no_show_rate IS 'Expected no-show rate';
COMMENT ON COLUMN overbooking_rules.seasonal_factor IS 'Seasonal adjustment factor';
COMMENT ON COLUMN overbooking_rules.effective_to IS 'NULL means indefinite';

-- Create overbooking_history table
CREATE TABLE overbooking_history (
    id BIGSERIAL PRIMARY KEY,
    schedule_id BIGINT NOT NULL REFERENCES flight_schedules(id),
    class_type class_type NOT NULL,
    total_bookings INTEGER NOT NULL,
    actual_capacity INTEGER NOT NULL,
    overbooking_count INTEGER NOT NULL,
    no_show_count INTEGER DEFAULT 0,
    denied_boarding_count INTEGER DEFAULT 0,
    compensation_amount DECIMAL(10,2) DEFAULT 0,
    resolution_status resolution_status DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for overbooking_history
CREATE INDEX idx_overbooking_history_schedule_id ON overbooking_history(schedule_id);
CREATE INDEX idx_overbooking_history_class_type ON overbooking_history(class_type);
CREATE INDEX idx_overbooking_history_created_at ON overbooking_history(created_at);
CREATE INDEX idx_overbooking_history_resolution_status ON overbooking_history(resolution_status);

-- Add comments
COMMENT ON TABLE overbooking_history IS 'Historical overbooking incidents';
COMMENT ON COLUMN overbooking_history.total_bookings IS 'Total confirmed bookings';
COMMENT ON COLUMN overbooking_history.actual_capacity IS 'Actual aircraft capacity';
COMMENT ON COLUMN overbooking_history.overbooking_count IS 'Number of overbooked passengers';
COMMENT ON COLUMN overbooking_history.no_show_count IS 'Actual no-show passengers';
COMMENT ON COLUMN overbooking_history.denied_boarding_count IS 'Passengers denied boarding';
COMMENT ON COLUMN overbooking_history.compensation_amount IS 'Total compensation paid';

-- Create compensation_rules table
CREATE TABLE compensation_rules (
    id SERIAL PRIMARY KEY,
    airline_id INTEGER NOT NULL REFERENCES airlines(id),
    region VARCHAR(10) NOT NULL,
    flight_type flight_type NOT NULL,
    delay_threshold_minutes INTEGER NOT NULL,
    compensation_amount DECIMAL(10,2) NOT NULL,
    compensation_type compensation_type NOT NULL,
    priority_class_multiplier DECIMAL(3,2) DEFAULT 1.00,
    active BOOLEAN DEFAULT TRUE,
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for compensation_rules updated_at
CREATE TRIGGER update_compensation_rules_updated_at 
    BEFORE UPDATE ON compensation_rules 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for compensation_rules
CREATE INDEX idx_compensation_rules_airline_id ON compensation_rules(airline_id);
CREATE INDEX idx_compensation_rules_region ON compensation_rules(region);
CREATE INDEX idx_compensation_rules_effective ON compensation_rules(effective_from, effective_to);
CREATE INDEX idx_compensation_rules_active ON compensation_rules(active);

-- Add comments
COMMENT ON TABLE compensation_rules IS 'Passenger compensation rules';
COMMENT ON COLUMN compensation_rules.region IS 'Regulatory region (EU, US, etc.)';
COMMENT ON COLUMN compensation_rules.delay_threshold_minutes IS 'Minimum delay for compensation';
COMMENT ON COLUMN compensation_rules.compensation_amount IS 'Compensation amount';
COMMENT ON COLUMN compensation_rules.priority_class_multiplier IS 'Multiplier for higher classes';
COMMENT ON COLUMN compensation_rules.effective_to IS 'NULL means indefinite';