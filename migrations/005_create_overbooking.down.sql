-- Drop overbooking-related tables in reverse order
DROP TABLE IF EXISTS compensation_rules CASCADE;
DROP TABLE IF EXISTS overbooking_history CASCADE;
DROP TABLE IF EXISTS overbooking_rules CASCADE;