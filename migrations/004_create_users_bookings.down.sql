-- Drop booking-related tables in reverse order
DROP TABLE IF EXISTS booking_status_log CASCADE;
DROP TABLE IF EXISTS booking_payments CASCADE;
DROP TABLE IF EXISTS booking_passengers CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS users CASCADE;