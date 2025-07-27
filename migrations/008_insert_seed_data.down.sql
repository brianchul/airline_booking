-- Remove seed data in reverse order of dependencies
DELETE FROM popular_routes;
DELETE FROM system_settings;
DELETE FROM compensation_rules;
DELETE FROM overbooking_rules;
DELETE FROM flight_inventory;
DELETE FROM flight_schedules;
DELETE FROM flights;
DELETE FROM aircraft;
DELETE FROM airports;
DELETE FROM airlines;
DELETE FROM users;