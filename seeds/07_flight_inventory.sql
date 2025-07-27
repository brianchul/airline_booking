-- Flight inventory for all scheduled flights
-- This creates realistic seat availability

DO $$
DECLARE
    schedule_rec RECORD;
    aircraft_rec RECORD;
    availability_factor DECIMAL;
    base_price_economy DECIMAL;
    base_price_business DECIMAL;
    base_price_first DECIMAL;
BEGIN
    -- Loop through each flight schedule
    FOR schedule_rec IN 
        SELECT fs.id as schedule_id, fs.flight_id, fs.departure_time, f.base_price_economy, f.base_price_business, f.base_price_first
        FROM flight_schedules fs
        JOIN flights f ON fs.flight_id = f.id
        WHERE fs.departure_time > CURRENT_TIMESTAMP
        ORDER BY fs.id
    LOOP
        -- Get aircraft configuration for this flight
        SELECT a.capacity_economy, a.capacity_business, a.capacity_first
        INTO aircraft_rec
        FROM flights f
        JOIN aircraft a ON f.aircraft_id = a.id
        WHERE f.id = schedule_rec.flight_id;
        
        -- Calculate availability factor (0.6 to 0.95 - some flights are more booked)
        availability_factor := 0.6 + (RANDOM() * 0.35);
        
        -- Dynamic pricing based on departure time and availability
        base_price_economy := schedule_rec.base_price_economy * (0.8 + (RANDOM() * 0.4));
        base_price_business := COALESCE(schedule_rec.base_price_business * (0.8 + (RANDOM() * 0.4)), 0);
        base_price_first := COALESCE(schedule_rec.base_price_first * (0.8 + (RANDOM() * 0.4)), 0);
        
        -- Economy class inventory
        IF aircraft_rec.capacity_economy > 0 THEN
            INSERT INTO flight_inventory (
                schedule_id, class_type, total_seats, available_seats, 
                current_price, overbooking_limit
            ) VALUES (
                schedule_rec.schedule_id,
                'ECONOMY',
                aircraft_rec.capacity_economy,
                FLOOR(aircraft_rec.capacity_economy * availability_factor),
                base_price_economy,
                FLOOR(aircraft_rec.capacity_economy * 0.08)  -- 8% overbooking
            );
        END IF;
        
        -- Business class inventory
        IF aircraft_rec.capacity_business > 0 THEN
            INSERT INTO flight_inventory (
                schedule_id, class_type, total_seats, available_seats, 
                current_price, overbooking_limit
            ) VALUES (
                schedule_rec.schedule_id,
                'BUSINESS',
                aircraft_rec.capacity_business,
                FLOOR(aircraft_rec.capacity_business * (availability_factor + 0.1)),  -- Business less booked
                base_price_business,
                FLOOR(aircraft_rec.capacity_business * 0.05)  -- 5% overbooking
            );
        END IF;
        
        -- First class inventory
        IF aircraft_rec.capacity_first > 0 THEN
            INSERT INTO flight_inventory (
                schedule_id, class_type, total_seats, available_seats, 
                current_price, overbooking_limit
            ) VALUES (
                schedule_rec.schedule_id,
                'FIRST',
                aircraft_rec.capacity_first,
                FLOOR(aircraft_rec.capacity_first * (availability_factor + 0.2)),  -- First class less booked
                base_price_first,
                FLOOR(aircraft_rec.capacity_first * 0.02)  -- 2% overbooking
            );
        END IF;
        
    END LOOP;
    
    -- Create some sold out flights for demo
    UPDATE flight_inventory 
    SET available_seats = 0, confirmed_seats = total_seats
    WHERE schedule_id % 31 = 0 AND class_type = 'ECONOMY';
    
    -- Create some nearly sold out flights
    UPDATE flight_inventory 
    SET available_seats = 1 + (id % 3)
    WHERE schedule_id % 17 = 0 AND available_seats > 5;
    
END $$;