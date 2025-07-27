-- Create flight_inventory table for real-time seat availability
CREATE TABLE flight_inventory (
    id BIGSERIAL PRIMARY KEY,
    schedule_id BIGINT NOT NULL REFERENCES flight_schedules(id),
    class_type class_type NOT NULL,
    total_seats INTEGER NOT NULL,
    available_seats INTEGER NOT NULL,
    reserved_seats INTEGER DEFAULT 0,
    confirmed_seats INTEGER DEFAULT 0,
    blocked_seats INTEGER DEFAULT 0,
    overbooking_limit INTEGER DEFAULT 0,
    current_price DECIMAL(10,2) NOT NULL,
    version INTEGER DEFAULT 1,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraint to ensure seat consistency
    CONSTRAINT chk_inventory_seats CHECK (
        available_seats + reserved_seats + confirmed_seats + blocked_seats <= total_seats + overbooking_limit
    ),
    
    -- Unique constraint for schedule and class combination
    CONSTRAINT uk_schedule_class UNIQUE (schedule_id, class_type)
);

-- Create trigger for flight_inventory last_updated
CREATE TRIGGER update_flight_inventory_last_updated 
    BEFORE UPDATE ON flight_inventory 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Update the trigger to update last_updated instead of updated_at
DROP TRIGGER update_flight_inventory_last_updated ON flight_inventory;
CREATE OR REPLACE FUNCTION update_last_updated_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_flight_inventory_last_updated 
    BEFORE UPDATE ON flight_inventory 
    FOR EACH ROW 
    EXECUTE FUNCTION update_last_updated_column();

-- Create indexes for flight_inventory
CREATE INDEX idx_flight_inventory_schedule_id ON flight_inventory(schedule_id);
CREATE INDEX idx_flight_inventory_class_type ON flight_inventory(class_type);
CREATE INDEX idx_flight_inventory_availability ON flight_inventory(schedule_id, class_type, available_seats);
CREATE INDEX idx_flight_inventory_version ON flight_inventory(schedule_id, version);

-- Add comments
COMMENT ON TABLE flight_inventory IS 'Real-time flight seat inventory';
COMMENT ON COLUMN flight_inventory.total_seats IS 'Total seats for this class';
COMMENT ON COLUMN flight_inventory.available_seats IS 'Currently available seats';
COMMENT ON COLUMN flight_inventory.reserved_seats IS 'Temporarily reserved seats';
COMMENT ON COLUMN flight_inventory.confirmed_seats IS 'Confirmed booked seats';
COMMENT ON COLUMN flight_inventory.blocked_seats IS 'Blocked/maintenance seats';
COMMENT ON COLUMN flight_inventory.overbooking_limit IS 'Allowed overbooking count';
COMMENT ON COLUMN flight_inventory.current_price IS 'Current dynamic price';
COMMENT ON COLUMN flight_inventory.version IS 'Optimistic locking version';