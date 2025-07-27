-- Drop tables in reverse order
DROP TABLE IF EXISTS aircraft CASCADE;
DROP TABLE IF EXISTS airports CASCADE;
DROP TABLE IF EXISTS airlines CASCADE;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop custom types
DROP TYPE IF EXISTS setting_type;
DROP TYPE IF EXISTS resolution_status;
DROP TYPE IF EXISTS seat_status;
DROP TYPE IF EXISTS seat_type;
DROP TYPE IF EXISTS flight_type;
DROP TYPE IF EXISTS compensation_type;
DROP TYPE IF EXISTS overbooking_class_type;
DROP TYPE IF EXISTS route_type;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS class_type;
DROP TYPE IF EXISTS flight_status;