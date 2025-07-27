-- Sample bookings for demo (various statuses)
DO $$
DECLARE
    user_rec RECORD;
    schedule_rec RECORD;
    booking_uuid UUID;
    booking_id BIGINT;
    counter INTEGER := 0;
BEGIN
    -- Create sample bookings for demo
    FOR user_rec IN 
        SELECT id as user_id FROM users WHERE username != 'admin' LIMIT 3
    LOOP
        -- Create 2-3 bookings per user
        FOR schedule_rec IN 
            SELECT fs.id as schedule_id, fi.class_type, fi.current_price
            FROM flight_schedules fs
            JOIN flight_inventory fi ON fs.id = fi.schedule_id
            WHERE fs.departure_time > CURRENT_TIMESTAMP + INTERVAL '1 day'
            AND fi.available_seats > 0
            ORDER BY RANDOM()
            LIMIT 2 + (counter % 2)
        LOOP
            booking_uuid := gen_random_uuid();
            
            INSERT INTO bookings (
                booking_uuid, user_id, schedule_id, class_type, 
                passenger_count, total_amount, status, booking_source,
                expires_at, seat_numbers
            ) VALUES (
                booking_uuid,
                user_rec.user_id,
                schedule_rec.schedule_id,
                schedule_rec.class_type,
                1 + (counter % 2),  -- 1 or 2 passengers
                schedule_rec.current_price * (1 + (counter % 2)),
                CASE 
                    WHEN counter % 5 = 0 THEN 'CONFIRMED'
                    WHEN counter % 5 = 1 THEN 'RESERVED'
                    WHEN counter % 5 = 2 THEN 'EXPIRED'
                    WHEN counter % 5 = 3 THEN 'CANCELLED'
                    ELSE 'PENDING'
                END,
                CASE WHEN counter % 3 = 0 THEN 'MOBILE' ELSE 'WEB' END,
                CASE 
                    WHEN counter % 5 = 1 THEN CURRENT_TIMESTAMP + INTERVAL '15 minutes'
                    ELSE NULL
                END,
                CASE 
                    WHEN counter % 5 = 0 THEN '["12A", "12B"]'::jsonb
                    WHEN counter % 5 = 1 THEN '["15C"]'::jsonb
                    ELSE NULL
                END
            ) RETURNING id INTO booking_id;
            
            -- Add passenger details
            INSERT INTO booking_passengers (
                booking_id, first_name, last_name, date_of_birth,
                passport_number, nationality, seat_number
            ) VALUES (
                booking_id,
                CASE counter % 4
                    WHEN 0 THEN 'Alice'
                    WHEN 1 THEN 'Bob'
                    WHEN 2 THEN 'Charlie'
                    ELSE 'Diana'
                END,
                CASE counter % 3
                    WHEN 0 THEN 'Brown'
                    WHEN 1 THEN 'Green'
                    ELSE 'White'
                END,
                '1990-01-01'::date + (counter * 100 || ' days')::interval,
                'P' || LPAD((1000000 + counter)::text, 8, '0'),
                'US',
                CASE 
                    WHEN counter % 5 = 0 THEN '12A'
                    WHEN counter % 5 = 1 THEN '15C'
                    ELSE NULL
                END
            );
            
            -- Add payment records for confirmed bookings
            IF counter % 5 = 0 THEN  -- CONFIRMED bookings
                INSERT INTO booking_payments (
                    booking_id, payment_method, amount, currency,
                    payment_status, transaction_id, processed_at
                ) VALUES (
                    booking_id,
                    CASE counter % 3
                        WHEN 0 THEN 'CREDIT_CARD'
                        WHEN 1 THEN 'DEBIT_CARD'
                        ELSE 'BANK_TRANSFER'
                    END,
                    schedule_rec.current_price * (1 + (counter % 2)),
                    'USD',
                    'COMPLETED',
                    'TXN_' || booking_uuid::text,
                    CURRENT_TIMESTAMP - (counter || ' hours')::interval
                );
                
                UPDATE bookings 
                SET confirmed_at = CURRENT_TIMESTAMP - (counter || ' hours')::interval
                WHERE id = booking_id;
            END IF;
            
            -- Add status log
            INSERT INTO booking_status_log (
                booking_id, old_status, new_status, reason, changed_by
            ) VALUES (
                booking_id,
                NULL,
                CASE 
                    WHEN counter % 5 = 0 THEN 'CONFIRMED'
                    WHEN counter % 5 = 1 THEN 'RESERVED'
                    WHEN counter % 5 = 2 THEN 'EXPIRED'
                    WHEN counter % 5 = 3 THEN 'CANCELLED'
                    ELSE 'PENDING'
                END,
                'Initial booking status',
                'SYSTEM'
            );
            
            counter := counter + 1;
        END LOOP;
    END LOOP;
END $$;