-- Demo users for JWT authentication testing
-- Password for all users: "password123" (hashed with bcrypt)

INSERT INTO users (username, password_hash, email, first_name, last_name, phone) VALUES
('admin', '$2a$10$v9RZSrtXi5rwgMzDpYHM7ujqvswgStgE1W4wG5KHe1pg6ejC.Bek2', 'admin@airline.com', 'Admin', 'User', '+1-555-0001'),
('john.doe', '$2a$10$v9RZSrtXi5rwgMzDpYHM7ujqvswgStgE1W4wG5KHe1pg6ejC.Bek2', 'john.doe@example.com', 'John', 'Doe', '+1-555-0002'),
('jane.smith', '$2a$10$v9RZSrtXi5rwgMzDpYHM7ujqvswgStgE1W4wG5KHe1pg6ejC.Bek2', 'jane.smith@example.com', 'Jane', 'Smith', '+1-555-0003'),
('mike.wilson', '$2a$10$v9RZSrtXi5rwgMzDpYHM7ujqvswgStgE1W4wG5KHe1pg6ejC.Bek2', 'mike.wilson@example.com', 'Mike', 'Wilson', '+1-555-0004'),
('sarah.johnson', '$2a$10$v9RZSrtXi5rwgMzDpYHM7ujqvswgStgE1W4wG5KHe1pg6ejC.Bek2', 'sarah.johnson@example.com', 'Sarah', 'Johnson', '+1-555-0005'),
('demo.user', '$2a$10$v9RZSrtXi5rwgMzDpYHM7ujqvswgStgE1W4wG5KHe1pg6ejC.Bek2', 'demo@example.com', 'Demo', 'User', '+1-555-0006')
ON CONFLICT (username) DO NOTHING;