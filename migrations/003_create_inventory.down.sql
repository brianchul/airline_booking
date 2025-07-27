-- Drop inventory table and related functions
DROP TABLE IF EXISTS flight_inventory CASCADE;
DROP FUNCTION IF EXISTS update_last_updated_column();