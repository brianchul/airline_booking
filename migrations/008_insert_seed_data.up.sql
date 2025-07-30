-- Insert seed data for testing and JWT authentication

-- Insert test users (3 hardcoded users for JWT authentication)
INSERT INTO users (username, password_hash, email, first_name, last_name, tier) VALUES
('admin', '$2a$10$v9RZSrtXi5rwgMzDpYHM7ujqvswgStgE1W4wG5KHe1pg6ejC.Bek2', 'admin@airline.com', 'Admin', 'User', 'vip'),
('user1', '$2a$10$v9RZSrtXi5rwgMzDpYHM7ujqvswgStgE1W4wG5KHe1pg6ejC.Bek2', 'user1@example.com', 'John', 'Doe', 'normal'),
('user2', '$2a$10$v9RZSrtXi5rwgMzDpYHM7ujqvswgStgE1W4wG5KHe1pg6ejC.Bek2', 'user2@example.com', 'Jane', 'Smith', 'vip');

-- Insert sample airlines
INSERT INTO airlines (code, name, country) VALUES
('AA', 'American Airlines', 'US'),
('UA', 'United Airlines', 'US'),
('DL', 'Delta Air Lines', 'US'),
('BA', 'British Airways', 'GB'),
('LH', 'Lufthansa', 'DE'),
('AF', 'Air France', 'FR'),
('JL', 'Japan Airlines', 'JP'),
('SQ', 'Singapore Airlines', 'SG');

-- Insert sample airports
INSERT INTO airports (code, name, city, country, timezone) VALUES
('JFK', 'John F. Kennedy International Airport', 'New York', 'US', 'America/New_York'),
('LAX', 'Los Angeles International Airport', 'Los Angeles', 'US', 'America/Los_Angeles'),
('ORD', 'Chicago O''Hare International Airport', 'Chicago', 'US', 'America/Chicago'),
('LHR', 'London Heathrow Airport', 'London', 'GB', 'Europe/London'),
('CDG', 'Charles de Gaulle Airport', 'Paris', 'FR', 'Europe/Paris'),
('FRA', 'Frankfurt Airport', 'Frankfurt', 'DE', 'Europe/Berlin'),
('NRT', 'Narita International Airport', 'Tokyo', 'JP', 'Asia/Tokyo'),
('SIN', 'Singapore Changi Airport', 'Singapore', 'SG', 'Asia/Singapore'),
('DXB', 'Dubai International Airport', 'Dubai', 'AE', 'Asia/Dubai'),
('SYD', 'Sydney Kingsford Smith Airport', 'Sydney', 'AU', 'Australia/Sydney');

-- Insert sample aircraft
INSERT INTO aircraft (model, manufacturer, capacity_economy, capacity_business, capacity_first) VALUES
('Boeing 737-800', 'Boeing', 162, 16, 0),
('Boeing 777-300ER', 'Boeing', 296, 42, 8),
('Airbus A320', 'Airbus', 150, 12, 0),
('Airbus A380', 'Airbus', 469, 78, 14),
('Boeing 787-9', 'Boeing', 254, 28, 6),
('Airbus A350-900', 'Airbus', 276, 42, 6);

-- Insert sample flights
INSERT INTO flights (flight_number, airline_id, departure_airport_id, arrival_airport_id, aircraft_id, 
                    base_price_economy, base_price_business, base_price_first, duration_minutes, distance_km) VALUES
('AA100', 1, 1, 2, 2, 450.00, 1200.00, 2500.00, 360, 3944),  -- JFK to LAX
('UA200', 2, 3, 1, 1, 280.00, 800.00, NULL, 150, 1145),      -- ORD to JFK
('DL300', 3, 2, 4, 5, 850.00, 2200.00, 4500.00, 660, 8750), -- LAX to LHR
('BA400', 4, 4, 5, 3, 320.00, 950.00, NULL, 80, 344),        -- LHR to CDG
('LH500', 5, 6, 7, 2, 1200.00, 3500.00, 7000.00, 695, 9560), -- FRA to NRT
('SQ600', 8, 8, 9, 4, 680.00, 1800.00, 3600.00, 460, 5320); -- SIN to DXB

-- Insert sample flight schedules for the next few days
INSERT INTO flight_schedules (flight_id, departure_time, arrival_time) VALUES
(1, '2024-07-28 08:00:00', '2024-07-28 14:00:00'),
(1, '2024-07-29 08:00:00', '2024-07-29 14:00:00'),
(1, '2024-07-30 08:00:00', '2024-07-30 14:00:00'),
(2, '2024-07-28 10:30:00', '2024-07-28 13:00:00'),
(2, '2024-07-29 10:30:00', '2024-07-29 13:00:00'),
(3, '2024-07-28 18:00:00', '2024-07-29 12:00:00'),
(3, '2024-07-30 18:00:00', '2024-07-31 12:00:00'),
(4, '2024-07-28 14:20:00', '2024-07-28 15:40:00'),
(5, '2024-07-29 22:30:00', '2024-07-30 20:05:00'),
(6, '2024-07-28 16:45:00', '2024-07-28 23:25:00');

-- Insert flight inventory for the scheduled flights
INSERT INTO flight_inventory (schedule_id, class_type, total_seats, available_seats, current_price) VALUES
-- Flight 1 schedules (Boeing 777-300ER: 296 economy, 42 business, 8 first)
(1, 'ECONOMY', 296, 280, 450.00),
(1, 'BUSINESS', 42, 35, 1200.00),
(1, 'FIRST', 8, 6, 2500.00),
(2, 'ECONOMY', 296, 275, 460.00),
(2, 'BUSINESS', 42, 38, 1250.00),
(2, 'FIRST', 8, 7, 2500.00),
(3, 'ECONOMY', 296, 285, 470.00),
(3, 'BUSINESS', 42, 40, 1180.00),
(3, 'FIRST', 8, 8, 2500.00),

-- Flight 2 schedules (Boeing 737-800: 162 economy, 16 business)
(4, 'ECONOMY', 162, 150, 280.00),
(4, 'BUSINESS', 16, 12, 800.00),
(5, 'ECONOMY', 162, 145, 290.00),
(5, 'BUSINESS', 16, 14, 850.00),

-- Flight 3 schedules (Boeing 787-9: 254 economy, 28 business, 6 first)
(6, 'ECONOMY', 254, 240, 850.00),
(6, 'BUSINESS', 28, 25, 2200.00),
(6, 'FIRST', 6, 5, 4500.00),
(7, 'ECONOMY', 254, 235, 880.00),
(7, 'BUSINESS', 28, 22, 2350.00),
(7, 'FIRST', 6, 4, 4700.00),

-- Flight 4 schedules (Airbus A320: 150 economy, 12 business)
(8, 'ECONOMY', 150, 140, 320.00),
(8, 'BUSINESS', 12, 10, 950.00),

-- Flight 5 schedules (Boeing 777-300ER)
(9, 'ECONOMY', 296, 250, 1200.00),
(9, 'BUSINESS', 42, 30, 3500.00),
(9, 'FIRST', 8, 6, 7000.00),

-- Flight 6 schedules (Airbus A380: 469 economy, 78 business, 14 first)
(10, 'ECONOMY', 469, 420, 680.00),
(10, 'BUSINESS', 78, 65, 1800.00),
(10, 'FIRST', 14, 12, 3600.00);

-- Insert sample overbooking rules
INSERT INTO overbooking_rules (airline_id, route_type, class_type, base_overbooking_rate, 
                              max_overbooking_rate, no_show_rate, effective_from) VALUES
(1, 'DOMESTIC', 'ECONOMY', 8.0, 12.0, 6.0, '2024-01-01'),
(1, 'INTERNATIONAL', 'ECONOMY', 5.0, 8.0, 4.0, '2024-01-01'),
(2, 'ALL', 'ECONOMY', 7.0, 10.0, 5.5, '2024-01-01'),
(3, 'ALL', 'ECONOMY', 6.5, 9.0, 5.0, '2024-01-01');

-- Insert sample compensation rules
INSERT INTO compensation_rules (airline_id, region, flight_type, delay_threshold_minutes, 
                               compensation_amount, compensation_type, effective_from) VALUES
(1, 'US', 'DOMESTIC', 180, 400.00, 'CASH', '2024-01-01'),
(1, 'US', 'INTERNATIONAL', 240, 600.00, 'CASH', '2024-01-01'),
(4, 'EU', 'DOMESTIC', 120, 250.00, 'CASH', '2024-01-01'),
(4, 'EU', 'INTERNATIONAL', 180, 400.00, 'CASH', '2024-01-01');

-- Insert system settings
INSERT INTO system_settings (setting_key, setting_value, setting_type, description, is_public) VALUES
('booking_expiry_minutes', '15', 'INTEGER', 'Booking reservation expiry time in minutes', TRUE),
('max_passengers_per_booking', '9', 'INTEGER', 'Maximum passengers allowed per booking', TRUE),
('search_cache_ttl_seconds', '300', 'INTEGER', 'Search results cache TTL in seconds', FALSE),
('inventory_cache_ttl_seconds', '120', 'INTEGER', 'Inventory cache TTL in seconds', FALSE),
('rate_limit_search_per_minute', '20', 'INTEGER', 'Search API rate limit per minute', FALSE),
('rate_limit_booking_per_minute', '5', 'INTEGER', 'Booking API rate limit per minute', FALSE);

-- Insert popular routes
INSERT INTO popular_routes (departure_airport_id, arrival_airport_id, search_count, booking_count, score) VALUES
(1, 2, 1250, 89, 95.5),  -- JFK to LAX
(2, 1, 1180, 82, 92.3),  -- LAX to JFK
(3, 1, 890, 67, 88.1),   -- ORD to JFK
(1, 4, 756, 45, 78.9),   -- JFK to LHR
(4, 1, 723, 42, 76.2);   -- LHR to JFK