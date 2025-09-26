-- +goose Up
CREATE TABLE workouts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    duration_minutes INTEGER NOT NULL,
    planned_date TIMESTAMP NOT NULL,
    date_completed TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose statementbegin
INSERT INTO workouts (
    id,
    user_id,
    title,
    description,
    duration_minutes,
    planned_date,
    date_completed,
    created_at,
    updated_at
) VALUES
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    'upper body strength',
    'bench press, shoulder press, pull-ups, biceps curls, triceps dips.',
    60,
    now() - INTERVAL '14 days',
    now() - INTERVAL '14 days',
    now() - INTERVAL '14 days',
    now() - INTERVAL '14 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    'leg day',
    'squats, lunges, deadlifts, calf raises, hip thrusts.',
    70,
    now() - INTERVAL '10 days',
    now() - INTERVAL '10 days',
    now() - INTERVAL '10 days',
    now() - INTERVAL '10 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    'cardio intervals',
    'treadmill intervals alternating between sprint and jog.',
    45,
    now() - INTERVAL '5 days',
    null,
    now() - INTERVAL '5 days',
    now() - INTERVAL '5 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    'full body circuit',
    'mix of strength and cardio: burpees, kettlebell swings, push-ups, '
    || 'rows, mountain climbers.',
    55,
    now() + INTERVAL '2 days',
    null,
    now(),
    now()
);
-- +goose statementend

-- +goose Down
DROP TABLE workouts;
