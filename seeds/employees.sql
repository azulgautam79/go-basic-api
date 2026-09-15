-- Seed data for employees table
INSERT INTO employees (id, first_name, last_name, email, gender, is_active, created_at, updated_at, deleted_at)
VALUES 
    (
        'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 
        'John', 
        'Doe', 
        'john.doe@example.com', 
        'male', 
        TRUE, 
        NOW() - INTERVAL '30 days', 
        NOW() - INTERVAL '30 days', 
        NULL
    ),
    (
        'b2c3d4e5-f6a7-8901-bcde-f12345678901', 
        'Jane', 
        'Smith', 
        'jane.smith@example.com', 
        'female', 
        TRUE, 
        NOW() - INTERVAL '25 days', 
        NOW() - INTERVAL '25 days', 
        NULL
    ),
    (
        'c3d4e5f6-a7b8-9012-cdef-123456789012', 
        'Alex', 
        'Taylor', 
        'alex.taylor@example.com', 
        'other', 
        TRUE, 
        NOW() - INTERVAL '20 days', 
        NOW() - INTERVAL '20 days', 
        NULL
    ),
    (
        'd4e5f6a7-b8c9-0123-def1-234567890123', 
        'Morgan', 
        'Lee', 
        'morgan.lee@example.com', 
        'prefer_not_to_say', 
        FALSE, 
        NOW() - INTERVAL '15 days', 
        NOW() - INTERVAL '15 days', 
        NULL
    ),
    (
        'e5f6a7b8-c9d0-1234-ef12-345678901234', 
        'Michael', 
        'Brown', 
        'michael.brown@example.com', 
        'male', 
        TRUE, 
        NOW() - INTERVAL '10 days', 
        NOW() - INTERVAL '10 days', 
        NULL
    ),
    (
        'f6a7b8c9-d0e1-2345-f123-456789012345', 
        'Emily', 
        'Davis', 
        'emily.davis@example.com', 
        'female', 
        FALSE, 
        NOW() - INTERVAL '5 days', 
        NOW() - INTERVAL '5 days', 
        NULL
    ),
    (
        'a0b1c2d3-e4f5-6789-0123-456789012345', 
        'Chris', 
        'Wilson', 
        'chris.wilson@example.com', 
        'male', 
        TRUE, 
        NOW() - INTERVAL '2 days', 
        NOW() - INTERVAL '2 days', 
        NOW() - INTERVAL '1 day' -- Soft deleted record
    )
ON CONFLICT (email) DO NOTHING;