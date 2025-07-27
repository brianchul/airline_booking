-- Flight schedules for the next 30 days (demo data)
-- This creates multiple daily flights for testing

DO $$
DECLARE
    flight_rec RECORD;
    schedule_date DATE;
    departure_base_time TIME;
    arrival_base_time TIME;
    flight_duration INTERVAL;
BEGIN
    -- Loop through each flight
    FOR flight_rec IN 
        SELECT f.id as flight_id, f.duration_minutes, f.flight_number
        FROM flights f 
        ORDER BY f.id
    LOOP
        -- Set base departure times based on flight number pattern
        CASE 
            WHEN flight_rec.flight_number LIKE '%100' THEN 
                departure_base_time := '08:00:00';
            WHEN flight_rec.flight_number LIKE '%200' THEN 
                departure_base_time := '14:30:00';
            WHEN flight_rec.flight_number LIKE '%300' THEN 
                departure_base_time := '20:15:00';
            ELSE 
                departure_base_time := '12:00:00';
        END CASE;
        
        flight_duration := (flight_rec.duration_minutes || ' minutes')::INTERVAL;
        arrival_base_time := departure_base_time + flight_duration;
        
        -- Create schedules for next 30 days
        FOR i IN 0..29 LOOP
            schedule_date := CURRENT_DATE + i;
            
            -- Skip some flights on certain days to create realistic availability
            CONTINUE WHEN (
                (EXTRACT(DOW FROM schedule_date) = 0 AND flight_rec.flight_number LIKE '%300') OR  -- No night flights on Sunday
                (i % 7 = 0 AND flight_rec.flight_number LIKE '%200')  -- Skip some afternoon flights weekly
            );
            
            INSERT INTO flight_schedules (
                flight_id, 
                departure_time, 
                arrival_time,
                status,
                gate,
                terminal
            ) VALUES (
                flight_rec.flight_id,
                schedule_date + departure_base_time,
                schedule_date + arrival_base_time,
                'SCHEDULED',
                'A' || (1 + (i % 20))::TEXT,  -- Gates A1-A20
                CASE WHEN (i % 3) = 0 THEN '1' ELSE '2' END  -- Terminal 1 or 2
            )
            ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
    
    -- Add some delayed and cancelled flights for realism
    UPDATE flight_schedules 
    SET status = 'DELAYED', delay_minutes = 45 
    WHERE id % 23 = 0 AND departure_time > CURRENT_TIMESTAMP;
    
    UPDATE flight_schedules 
    SET status = 'CANCELLED', cancellation_reason = 'Weather conditions' 
    WHERE id % 47 = 0 AND departure_time > CURRENT_TIMESTAMP;
    
END $$;