-- Aircraft types and configurations for demo
INSERT INTO aircraft (model, manufacturer, capacity_economy, capacity_business, capacity_first) VALUES
-- Boeing Aircraft
('Boeing 737-800', 'Boeing', 162, 16, 0),
('Boeing 737-900', 'Boeing', 178, 20, 0),
('Boeing 757-200', 'Boeing', 200, 16, 0),
('Boeing 767-300', 'Boeing', 218, 30, 0),
('Boeing 777-200', 'Boeing', 314, 37, 8),
('Boeing 777-300ER', 'Boeing', 296, 42, 8),
('Boeing 787-8', 'Boeing', 234, 28, 0),
('Boeing 787-9', 'Boeing', 254, 28, 6),
('Boeing 787-10', 'Boeing', 318, 28, 6),

-- Airbus Aircraft
('Airbus A319', 'Airbus', 134, 8, 0),
('Airbus A320', 'Airbus', 150, 12, 0),
('Airbus A321', 'Airbus', 185, 16, 0),
('Airbus A330-200', 'Airbus', 247, 28, 0),
('Airbus A330-300', 'Airbus', 295, 42, 0),
('Airbus A340-300', 'Airbus', 267, 30, 12),
('Airbus A350-900', 'Airbus', 276, 42, 6),
('Airbus A350-1000', 'Airbus', 327, 48, 9),
('Airbus A380', 'Airbus', 469, 78, 14),

-- Regional Aircraft
('Embraer E-Jet 190', 'Embraer', 96, 4, 0),
('Bombardier CRJ-900', 'Bombardier', 76, 4, 0),

-- Cargo Converted (for cargo flights)
('Boeing 747-8F', 'Boeing', 0, 0, 0),
('Airbus A330-200F', 'Airbus', 0, 0, 0)
ON CONFLICT (model, manufacturer) DO NOTHING;