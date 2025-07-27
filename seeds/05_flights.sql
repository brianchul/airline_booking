-- Popular flight routes for demo
INSERT INTO flights (flight_number, airline_id, departure_airport_id, arrival_airport_id, aircraft_id, 
                    base_price_economy, base_price_business, base_price_first, duration_minutes, distance_km) VALUES

-- American Airlines Routes
('AA100', (SELECT id FROM airlines WHERE code = 'AA'), (SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM airports WHERE code = 'LAX'), (SELECT id FROM aircraft WHERE model = 'Boeing 777-300ER' LIMIT 1), 450.00, 1200.00, 2500.00, 360, 3944),
('AA200', (SELECT id FROM airlines WHERE code = 'AA'), (SELECT id FROM airports WHERE code = 'LAX'), (SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM aircraft WHERE model = 'Boeing 777-300ER' LIMIT 1), 450.00, 1200.00, 2500.00, 330, 3944),
('AA300', (SELECT id FROM airlines WHERE code = 'AA'), (SELECT id FROM airports WHERE code = 'DFW'), (SELECT id FROM airports WHERE code = 'LHR'), (SELECT id FROM aircraft WHERE model = 'Boeing 787-9' LIMIT 1), 850.00, 2200.00, 4500.00, 540, 7800),

-- United Airlines Routes
('UA100', (SELECT id FROM airlines WHERE code = 'UA'), (SELECT id FROM airports WHERE code = 'ORD'), (SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM aircraft WHERE model = 'Boeing 737-800' LIMIT 1), 280.00, 800.00, NULL, 150, 1145),
('UA200', (SELECT id FROM airlines WHERE code = 'UA'), (SELECT id FROM airports WHERE code = 'SFO'), (SELECT id FROM airports WHERE code = 'NRT'), (SELECT id FROM aircraft WHERE model = 'Boeing 777-300ER' LIMIT 1), 950.00, 2800.00, 5500.00, 645, 8280),
('UA300', (SELECT id FROM airlines WHERE code = 'UA'), (SELECT id FROM airports WHERE code = 'DEN'), (SELECT id FROM airports WHERE code = 'FRA'), (SELECT id FROM aircraft WHERE model = 'Boeing 787-9' LIMIT 1), 780.00, 2100.00, 4200.00, 600, 7920),

-- Delta Air Lines Routes
('DL100', (SELECT id FROM airlines WHERE code = 'DL'), (SELECT id FROM airports WHERE code = 'ATL'), (SELECT id FROM airports WHERE code = 'CDG'), (SELECT id FROM aircraft WHERE model = 'Airbus A350-900' LIMIT 1), 720.00, 1950.00, 3800.00, 520, 7000),
('DL200', (SELECT id FROM airlines WHERE code = 'DL'), (SELECT id FROM airports WHERE code = 'LAX'), (SELECT id FROM airports WHERE code = 'LHR'), (SELECT id FROM aircraft WHERE model = 'Boeing 787-9' LIMIT 1), 850.00, 2200.00, 4500.00, 660, 8750),
('DL300', (SELECT id FROM airlines WHERE code = 'DL'), (SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM airports WHERE code = 'AMS'), (SELECT id FROM aircraft WHERE model = 'Airbus A330-300' LIMIT 1), 680.00, 1800.00, NULL, 460, 5850),

-- British Airways Routes
('BA100', (SELECT id FROM airlines WHERE code = 'BA'), (SELECT id FROM airports WHERE code = 'LHR'), (SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM aircraft WHERE model = 'Boeing 777-300ER' LIMIT 1), 750.00, 2100.00, 4200.00, 480, 5540),
('BA200', (SELECT id FROM airlines WHERE code = 'BA'), (SELECT id FROM airports WHERE code = 'LHR'), (SELECT id FROM airports WHERE code = 'CDG'), (SELECT id FROM aircraft WHERE model = 'Airbus A320' LIMIT 1), 180.00, 520.00, NULL, 80, 344),
('BA300', (SELECT id FROM airlines WHERE code = 'BA'), (SELECT id FROM airports WHERE code = 'LHR'), (SELECT id FROM airports WHERE code = 'SIN'), (SELECT id FROM aircraft WHERE model = 'Airbus A350-900' LIMIT 1), 1200.00, 3200.00, 6500.00, 780, 10890),

-- Lufthansa Routes
('LH100', (SELECT id FROM airlines WHERE code = 'LH'), (SELECT id FROM airports WHERE code = 'FRA'), (SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM aircraft WHERE model = 'Boeing 747-8F' LIMIT 1), 720.00, 2000.00, 4000.00, 500, 6194),
('LH200', (SELECT id FROM airlines WHERE code = 'LH'), (SELECT id FROM airports WHERE code = 'FRA'), (SELECT id FROM airports WHERE code = 'NRT'), (SELECT id FROM aircraft WHERE model = 'Boeing 777-300ER' LIMIT 1), 1100.00, 3200.00, 6500.00, 695, 9560),
('LH300', (SELECT id FROM airlines WHERE code = 'LH'), (SELECT id FROM airports WHERE code = 'FRA'), (SELECT id FROM airports WHERE code = 'SIN'), (SELECT id FROM aircraft WHERE model = 'Airbus A350-900' LIMIT 1), 950.00, 2800.00, 5600.00, 740, 10350),

-- Singapore Airlines Routes
('SQ100', (SELECT id FROM airlines WHERE code = 'SQ'), (SELECT id FROM airports WHERE code = 'SIN'), (SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM aircraft WHERE model = 'Airbus A350-1000' LIMIT 1), 1400.00, 3800.00, 7500.00, 1020, 17460),
('SQ200', (SELECT id FROM airlines WHERE code = 'SQ'), (SELECT id FROM airports WHERE code = 'SIN'), (SELECT id FROM airports WHERE code = 'LHR'), (SELECT id FROM aircraft WHERE model = 'Airbus A380' LIMIT 1), 1100.00, 3000.00, 6000.00, 780, 10890),
('SQ300', (SELECT id FROM airlines WHERE code = 'SQ'), (SELECT id FROM airports WHERE code = 'SIN'), (SELECT id FROM airports WHERE code = 'SYD'), (SELECT id FROM aircraft WHERE model = 'Boeing 777-300ER' LIMIT 1), 580.00, 1600.00, 3200.00, 480, 6300),

-- Emirates Routes
('EK100', (SELECT id FROM airlines WHERE code = 'EK'), (SELECT id FROM airports WHERE code = 'DXB'), (SELECT id FROM airports WHERE code = 'JFK'), (SELECT id FROM aircraft WHERE model = 'Airbus A380' LIMIT 1), 1200.00, 3500.00, 7000.00, 840, 11000),
('EK200', (SELECT id FROM airlines WHERE code = 'EK'), (SELECT id FROM airports WHERE code = 'DXB'), (SELECT id FROM airports WHERE code = 'LAX'), (SELECT id FROM aircraft WHERE model = 'Airbus A380' LIMIT 1), 1300.00, 3600.00, 7200.00, 900, 13420),
('EK300', (SELECT id FROM airlines WHERE code = 'EK'), (SELECT id FROM airports WHERE code = 'DXB'), (SELECT id FROM airports WHERE code = 'SYD'), (SELECT id FROM aircraft WHERE model = 'Airbus A380' LIMIT 1), 950.00, 2800.00, 5600.00, 850, 12050),

-- Japan Airlines Routes
('JL100', (SELECT id FROM airlines WHERE code = 'JL'), (SELECT id FROM airports WHERE code = 'NRT'), (SELECT id FROM airports WHERE code = 'LAX'), (SELECT id FROM aircraft WHERE model = 'Boeing 787-9' LIMIT 1), 850.00, 2400.00, 4800.00, 650, 8800),
('JL200', (SELECT id FROM airlines WHERE code = 'JL'), (SELECT id FROM airports WHERE code = 'HND'), (SELECT id FROM airports WHERE code = 'SIN'), (SELECT id FROM aircraft WHERE model = 'Boeing 787-8' LIMIT 1), 650.00, 1800.00, 3600.00, 420, 5300),

-- Qantas Routes
('QF100', (SELECT id FROM airlines WHERE code = 'QF'), (SELECT id FROM airports WHERE code = 'SYD'), (SELECT id FROM airports WHERE code = 'LAX'), (SELECT id FROM aircraft WHERE model = 'Airbus A380' LIMIT 1), 1100.00, 3200.00, 6400.00, 720, 12050),
('QF200', (SELECT id FROM airlines WHERE code = 'QF'), (SELECT id FROM airports WHERE code = 'SYD'), (SELECT id FROM airports WHERE code = 'LHR'), (SELECT id FROM aircraft WHERE model = 'Airbus A380' LIMIT 1), 1400.00, 4000.00, 8000.00, 1200, 17000)
ON CONFLICT (flight_number, airline_id) DO NOTHING;