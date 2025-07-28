-- Create popular_routes table
CREATE TABLE popular_routes (
    id SERIAL PRIMARY KEY,
    departure_airport_id INTEGER NOT NULL REFERENCES airports(id),
    arrival_airport_id INTEGER NOT NULL REFERENCES airports(id),
    search_count BIGINT DEFAULT 0,
    booking_count BIGINT DEFAULT 0,
    last_searched TIMESTAMP,
    last_booked TIMESTAMP,
    score DECIMAL(8,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint for route combination
    CONSTRAINT uk_route UNIQUE (departure_airport_id, arrival_airport_id)
);

-- Create trigger for popular_routes updated_at
CREATE TRIGGER update_popular_routes_updated_at 
    BEFORE UPDATE ON popular_routes 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for popular_routes
CREATE INDEX idx_popular_routes_score ON popular_routes(score);
CREATE INDEX idx_popular_routes_last_searched ON popular_routes(last_searched);

-- Add comments
COMMENT ON TABLE popular_routes IS 'Popular flight routes for cache optimization';
COMMENT ON COLUMN popular_routes.score IS 'Popularity score for cache prioritization';

-- Create search_analytics table without partitioning
CREATE TABLE search_analytics (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    departure_airport_id INTEGER REFERENCES airports(id),
    arrival_airport_id INTEGER REFERENCES airports(id),
    departure_date DATE,
    return_date DATE,
    passenger_count SMALLINT DEFAULT 1,
    class_preference class_type,
    search_result_count INTEGER,
    response_time_ms INTEGER,
    converted_to_booking BOOLEAN DEFAULT FALSE,
    user_agent TEXT,
    ip_address INET,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for search_analytics
CREATE INDEX idx_search_analytics_user_id ON search_analytics(user_id);
CREATE INDEX idx_search_analytics_route ON search_analytics(departure_airport_id, arrival_airport_id);
CREATE INDEX idx_search_analytics_departure_date ON search_analytics(departure_date);
CREATE INDEX idx_search_analytics_converted ON search_analytics(converted_to_booking);
CREATE INDEX idx_search_analytics_created_at ON search_analytics(created_at);

-- Add comments
COMMENT ON TABLE search_analytics IS 'Search behavior analytics';
COMMENT ON COLUMN search_analytics.user_id IS 'NULL for anonymous users';
COMMENT ON COLUMN search_analytics.return_date IS 'NULL for one-way';
COMMENT ON COLUMN search_analytics.search_result_count IS 'Number of results returned';
COMMENT ON COLUMN search_analytics.response_time_ms IS 'Search response time';