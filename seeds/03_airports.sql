-- Major airports worldwide for demo
INSERT INTO airports (code, name, city, country, timezone, latitude, longitude) VALUES
-- North American Airports
('JFK', 'John F. Kennedy International Airport', 'New York', 'US', 'America/New_York', 40.6413, -73.7781),
('LAX', 'Los Angeles International Airport', 'Los Angeles', 'US', 'America/Los_Angeles', 33.9425, -118.4081),
('ORD', 'Chicago O''Hare International Airport', 'Chicago', 'US', 'America/Chicago', 41.9742, -87.9073),
('ATL', 'Hartsfield-Jackson Atlanta International Airport', 'Atlanta', 'US', 'America/New_York', 33.6407, -84.4277),
('DFW', 'Dallas/Fort Worth International Airport', 'Dallas', 'US', 'America/Chicago', 32.8998, -97.0403),
('DEN', 'Denver International Airport', 'Denver', 'US', 'America/Denver', 39.8561, -104.6737),
('SFO', 'San Francisco International Airport', 'San Francisco', 'US', 'America/Los_Angeles', 37.6213, -122.3790),
('SEA', 'Seattle-Tacoma International Airport', 'Seattle', 'US', 'America/Los_Angeles', 47.4502, -122.3088),
('MIA', 'Miami International Airport', 'Miami', 'US', 'America/New_York', 25.7959, -80.2870),
('YYZ', 'Toronto Pearson International Airport', 'Toronto', 'CA', 'America/Toronto', 43.6777, -79.6248),

-- European Airports
('LHR', 'London Heathrow Airport', 'London', 'GB', 'Europe/London', 51.4700, -0.4543),
('CDG', 'Charles de Gaulle Airport', 'Paris', 'FR', 'Europe/Paris', 49.0097, 2.5479),
('FRA', 'Frankfurt Airport', 'Frankfurt', 'DE', 'Europe/Berlin', 50.0379, 8.5622),
('AMS', 'Amsterdam Airport Schiphol', 'Amsterdam', 'NL', 'Europe/Amsterdam', 52.3105, 4.7683),
('FCO', 'Leonardo da Vinci International Airport', 'Rome', 'IT', 'Europe/Rome', 41.8003, 12.2389),
('MAD', 'Adolfo Suárez Madrid-Barajas Airport', 'Madrid', 'ES', 'Europe/Madrid', 40.4936, -3.5668),
('ZUR', 'Zurich Airport', 'Zurich', 'CH', 'Europe/Zurich', 47.4647, 8.5492),

-- Asian Airports
('NRT', 'Narita International Airport', 'Tokyo', 'JP', 'Asia/Tokyo', 35.7720, 140.3928),
('HND', 'Haneda Airport', 'Tokyo', 'JP', 'Asia/Tokyo', 35.5494, 139.7798),
('ICN', 'Incheon International Airport', 'Seoul', 'KR', 'Asia/Seoul', 37.4602, 126.4407),
('SIN', 'Singapore Changi Airport', 'Singapore', 'SG', 'Asia/Singapore', 1.3644, 103.9915),
('HKG', 'Hong Kong International Airport', 'Hong Kong', 'HK', 'Asia/Hong_Kong', 22.3080, 113.9185),
('BKK', 'Suvarnabhumi Airport', 'Bangkok', 'TH', 'Asia/Bangkok', 13.6900, 100.7501),
('KUL', 'Kuala Lumpur International Airport', 'Kuala Lumpur', 'MY', 'Asia/Kuala_Lumpur', 2.7456, 101.7072),

-- Middle Eastern Airports
('DXB', 'Dubai International Airport', 'Dubai', 'AE', 'Asia/Dubai', 25.2532, 55.3657),
('DOH', 'Hamad International Airport', 'Doha', 'QA', 'Asia/Qatar', 25.2731, 51.6089),
('AUH', 'Abu Dhabi International Airport', 'Abu Dhabi', 'AE', 'Asia/Dubai', 24.4330, 54.6511),

-- Asia Pacific Airports
('SYD', 'Sydney Kingsford Smith Airport', 'Sydney', 'AU', 'Australia/Sydney', -33.9399, 151.1753),
('MEL', 'Melbourne Airport', 'Melbourne', 'AU', 'Australia/Melbourne', -37.6690, 144.8410),
('AKL', 'Auckland Airport', 'Auckland', 'NZ', 'Pacific/Auckland', -37.0082, 174.7850)
ON CONFLICT (code) DO NOTHING;