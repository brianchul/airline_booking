-- System configuration and rules for demo

-- Overbooking rules
INSERT INTO overbooking_rules (
    airline_id, aircraft_id, route_type, class_type, 
    base_overbooking_rate, max_overbooking_rate, no_show_rate,
    seasonal_factor, effective_from
) VALUES
-- American Airlines rules
((SELECT id FROM airlines WHERE code = 'AA'), NULL, 'DOMESTIC', 'ECONOMY', 8.0, 12.0, 6.0, 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'AA'), NULL, 'INTERNATIONAL', 'ECONOMY', 5.0, 8.0, 4.0, 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'AA'), NULL, 'ALL', 'BUSINESS', 3.0, 5.0, 2.0, 1.0, '2024-01-01'),

-- United Airlines rules
((SELECT id FROM airlines WHERE code = 'UA'), NULL, 'ALL', 'ECONOMY', 7.0, 10.0, 5.5, 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'UA'), NULL, 'ALL', 'BUSINESS', 2.5, 4.0, 1.5, 1.0, '2024-01-01'),

-- Delta Air Lines rules
((SELECT id FROM airlines WHERE code = 'DL'), NULL, 'ALL', 'ECONOMY', 6.5, 9.0, 5.0, 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'DL'), NULL, 'ALL', 'BUSINESS', 3.0, 5.0, 2.0, 1.0, '2024-01-01'),

-- European Airlines (more conservative)
((SELECT id FROM airlines WHERE code = 'BA'), NULL, 'ALL', 'ECONOMY', 4.0, 7.0, 3.5, 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'LH'), NULL, 'ALL', 'ECONOMY', 4.5, 7.5, 4.0, 1.0, '2024-01-01'),

-- Asian Airlines
((SELECT id FROM airlines WHERE code = 'SQ'), NULL, 'ALL', 'ECONOMY', 3.0, 6.0, 2.5, 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'JL'), NULL, 'ALL', 'ECONOMY', 3.5, 6.5, 3.0, 1.0, '2024-01-01')
ON CONFLICT DO NOTHING;

-- Compensation rules
INSERT INTO compensation_rules (
    airline_id, region, flight_type, delay_threshold_minutes,
    compensation_amount, compensation_type, priority_class_multiplier,
    effective_from
) VALUES
-- US regulations
((SELECT id FROM airlines WHERE code = 'AA'), 'US', 'DOMESTIC', 180, 400.00, 'CASH', 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'AA'), 'US', 'INTERNATIONAL', 240, 600.00, 'CASH', 1.5, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'UA'), 'US', 'DOMESTIC', 180, 400.00, 'CASH', 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'UA'), 'US', 'INTERNATIONAL', 240, 600.00, 'CASH', 1.5, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'DL'), 'US', 'DOMESTIC', 180, 400.00, 'CASH', 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'DL'), 'US', 'INTERNATIONAL', 240, 600.00, 'CASH', 1.5, '2024-01-01'),

-- EU regulations (EU261)
((SELECT id FROM airlines WHERE code = 'BA'), 'EU', 'DOMESTIC', 120, 250.00, 'CASH', 2.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'BA'), 'EU', 'INTERNATIONAL', 180, 400.00, 'CASH', 2.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'LH'), 'EU', 'DOMESTIC', 120, 250.00, 'CASH', 2.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'LH'), 'EU', 'INTERNATIONAL', 180, 400.00, 'CASH', 2.0, '2024-01-01'),

-- Asian regulations (more lenient)
((SELECT id FROM airlines WHERE code = 'SQ'), 'ASIA', 'DOMESTIC', 240, 200.00, 'VOUCHER', 1.5, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'SQ'), 'ASIA', 'INTERNATIONAL', 300, 300.00, 'VOUCHER', 1.5, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'JL'), 'ASIA', 'DOMESTIC', 240, 200.00, 'MILES', 1.0, '2024-01-01'),
((SELECT id FROM airlines WHERE code = 'JL'), 'ASIA', 'INTERNATIONAL', 300, 300.00, 'MILES', 1.0, '2024-01-01')
ON CONFLICT DO NOTHING;

-- System settings
INSERT INTO system_settings (setting_key, setting_value, setting_type, description, is_public) VALUES
-- Booking settings
('booking_expiry_minutes', '15', 'INTEGER', 'Booking reservation expiry time in minutes', TRUE),
('max_passengers_per_booking', '9', 'INTEGER', 'Maximum passengers allowed per booking', TRUE),
('min_booking_advance_hours', '2', 'INTEGER', 'Minimum hours before departure to allow booking', TRUE),
('max_booking_advance_days', '365', 'INTEGER', 'Maximum days in advance to allow booking', TRUE),

-- Cache settings
('search_cache_ttl_seconds', '300', 'INTEGER', 'Search results cache TTL in seconds', FALSE),
('inventory_cache_ttl_seconds', '120', 'INTEGER', 'Inventory cache TTL in seconds', FALSE),
('flight_details_cache_ttl_seconds', '600', 'INTEGER', 'Flight details cache TTL in seconds', FALSE),
('popular_routes_cache_ttl_seconds', '3600', 'INTEGER', 'Popular routes cache TTL in seconds', FALSE),

-- Rate limiting
('rate_limit_search_per_minute', '20', 'INTEGER', 'Search API rate limit per minute', FALSE),
('rate_limit_booking_per_minute', '5', 'INTEGER', 'Booking API rate limit per minute', FALSE),
('rate_limit_auth_per_minute', '10', 'INTEGER', 'Authentication API rate limit per minute', FALSE),

-- JWT settings
('jwt_secret_key', 'your-super-secret-jwt-key-change-in-production', 'STRING', 'JWT signing secret key', FALSE),
('jwt_expiry_hours', '24', 'INTEGER', 'JWT token expiry time in hours', FALSE),
('jwt_refresh_expiry_days', '7', 'INTEGER', 'JWT refresh token expiry in days', FALSE),

-- Business rules
('seat_hold_timeout_minutes', '10', 'INTEGER', 'How long to hold seats during selection', TRUE),
('price_change_threshold_percent', '15', 'INTEGER', 'Maximum price change percentage during booking', TRUE),
('inventory_update_interval_seconds', '30', 'INTEGER', 'How often to update inventory from cache', FALSE),

-- Notification settings
('email_enabled', 'true', 'BOOLEAN', 'Enable email notifications', FALSE),
('sms_enabled', 'true', 'BOOLEAN', 'Enable SMS notifications', FALSE),
('booking_confirmation_email', 'true', 'BOOLEAN', 'Send booking confirmation emails', TRUE),
('booking_reminder_hours', '24', 'INTEGER', 'Hours before departure to send reminder', TRUE),

-- Overbooking settings
('overbooking_enabled', 'true', 'BOOLEAN', 'Enable overbooking functionality', FALSE),
('auto_compensation_enabled', 'true', 'BOOLEAN', 'Enable automatic compensation calculation', FALSE),
('volunteer_compensation_multiplier', '1.5', 'STRING', 'Multiplier for volunteer compensation', FALSE),

-- API settings
('api_version', 'v1', 'STRING', 'Current API version', TRUE),
('maintenance_mode', 'false', 'BOOLEAN', 'System maintenance mode', TRUE),
('max_search_results', '100', 'INTEGER', 'Maximum search results per page', TRUE),
('default_page_size', '20', 'INTEGER', 'Default pagination page size', TRUE)
ON CONFLICT (setting_key) DO NOTHING;

-- Popular routes for cache optimization
INSERT INTO popular_routes (departure_airport_id, arrival_airport_id, search_count, booking_count, score, last_searched, last_booked) VALUES
((SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM airports WHERE code = 'LAX'), 2450, 189, 98.5, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '30 minutes'),
((SELECT id FROM airports WHERE code = 'LAX'), (SELECT id FROM airports WHERE code = 'JFK'), 2380, 182, 97.8, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '45 minutes'),
((SELECT id FROM airports WHERE code = 'ORD'), (SELECT id FROM airports WHERE code = 'JFK'), 1890, 167, 94.1, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
((SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM airports WHERE code = 'LHR'), 1756, 145, 91.9, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
((SELECT id FROM airports WHERE code = 'LHR'), (SELECT id FROM airports WHERE code = 'JFK'), 1723, 142, 91.2, CURRENT_TIMESTAMP - INTERVAL '4 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours'),
((SELECT id FROM airports WHERE code = 'SFO'), (SELECT id FROM airports WHERE code = 'NRT'), 1456, 98, 87.3, CURRENT_TIMESTAMP - INTERVAL '5 hours', CURRENT_TIMESTAMP - INTERVAL '4 hours'),
((SELECT id FROM airports WHERE code = 'DXB'), (SELECT id FROM airports WHERE code = 'JFK'), 1234, 78, 84.6, CURRENT_TIMESTAMP - INTERVAL '6 hours', CURRENT_TIMESTAMP - INTERVAL '5 hours'),
((SELECT id FROM airports WHERE code = 'SIN'), (SELECT id FROM airports WHERE code = 'LHR'), 1123, 67, 82.1, CURRENT_TIMESTAMP - INTERVAL '7 hours', CURRENT_TIMESTAMP - INTERVAL '6 hours'),
((SELECT id FROM airports WHERE code = 'FRA'), (SELECT id FROM airports WHERE code = 'JFK'), 1098, 65, 81.7, CURRENT_TIMESTAMP - INTERVAL '8 hours', CURRENT_TIMESTAMP - INTERVAL '7 hours'),
((SELECT id FROM airports WHERE code = 'CDG'), (SELECT id FROM airports WHERE code = 'JFK'), 987, 58, 78.9, CURRENT_TIMESTAMP - INTERVAL '9 hours', CURRENT_TIMESTAMP - INTERVAL '8 hours')
ON CONFLICT (departure_airport_id, arrival_airport_id) DO NOTHING;