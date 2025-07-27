-- Drop system tables in reverse order
DROP TABLE IF EXISTS api_rate_limits CASCADE;
DROP TABLE IF EXISTS system_settings CASCADE;
DROP TABLE IF EXISTS seat_assignments CASCADE;
DROP TABLE IF EXISTS seat_maps CASCADE;