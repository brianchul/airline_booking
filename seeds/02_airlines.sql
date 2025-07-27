-- Major airlines data for demo
INSERT INTO airlines (code, name, country) VALUES
-- North American Airlines
('AA', 'American Airlines', 'US'),
('UA', 'United Airlines', 'US'),
('DL', 'Delta Air Lines', 'US'),
('WN', 'Southwest Airlines', 'US'),
('AC', 'Air Canada', 'CA'),

-- European Airlines
('BA', 'British Airways', 'GB'),
('LH', 'Lufthansa', 'DE'),
('AF', 'Air France', 'FR'),
('KL', 'KLM Royal Dutch Airlines', 'NL'),
('AZ', 'Alitalia', 'IT'),
('IB', 'Iberia', 'ES'),

-- Asian Airlines
('JL', 'Japan Airlines', 'JP'),
('NH', 'All Nippon Airways', 'JP'),
('SQ', 'Singapore Airlines', 'SG'),
('CX', 'Cathay Pacific', 'HK'),
('TG', 'Thai Airways', 'TH'),
('KE', 'Korean Air', 'KR'),

-- Middle Eastern Airlines
('EK', 'Emirates', 'AE'),
('QR', 'Qatar Airways', 'QA'),
('ET', 'Etihad Airways', 'AE'),

-- Asia Pacific Airlines
('QF', 'Qantas', 'AU'),
('NZ', 'Air New Zealand', 'NZ')
ON CONFLICT (code) DO NOTHING;