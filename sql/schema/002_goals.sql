-- +goose Up
CREATE TABLE goals (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    goal_name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    goal_date TIMESTAMP NOT NULL,
    completion_date TIMESTAMP DEFAULT NULL,
    notes TEXT DEFAULT NULL,
    status VARCHAR(11) NOT NULL DEFAULT 'in_progress',
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (goal_name, user_id)
);

-- +goose statementbegin
INSERT INTO goals (
    id,
    created_at,
    updated_at,
    goal_name,
    description,
    goal_date,
    completion_date,
    notes,
    status,
    user_id
) VALUES
(
    gen_random_uuid(),
    now() - INTERVAL '30 days',
    now(),
    'run 5k without stopping',
    'train consistently to build endurance and complete a 5 kilometer run '
    || 'without walking breaks.',
    now() + INTERVAL '30 days',
    NULL,
    'currently able to run 3k comfortably. adding intervals to training.',
    'in_progress',
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    )
),
(
    gen_random_uuid(),
    now() - INTERVAL '90 days',
    now() - INTERVAL '60 days',
    'track meals daily for 30 days',
    'use a food journal or tracking app to log all meals, snacks, and drinks '
    || 'for accountability.',
    now() - INTERVAL '60 days',
    now() - INTERVAL '60 days',
    'learned a lot about portion sizes and calorie balance. built habit of '
    || 'awareness.',
    'completed',
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    )
),
(
    gen_random_uuid(),
    now() - INTERVAL '10 days',
    now(),
    'bench press 185 lbs',
    'increase strength progressively to achieve a one-rep max bench press of '
    || '185 pounds.',
    now() + INTERVAL '60 days',
    NULL,
    'currently at 165 lbs. adding 5 lbs per week with progressive overload.',
    'in_progress',
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    )
),
(
    gen_random_uuid(),
    now(),
    now(),
    'sleep 8 hours consistently',
    'establish a consistent sleep routine to achieve at least 8 hours of '
    || 'quality sleep per night.',
    now() + INTERVAL '90 days',
    NULL,
    'planning to set a bedtime alarm and reduce screen time before bed.',
    'in_progress',
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    )
);
-- +goose statementend

-- +goose Down
DROP TABLE goals;
