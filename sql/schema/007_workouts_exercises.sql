-- +goose Up
CREATE TABLE workouts_exercises (
    id UUID PRIMARY KEY,
    workout_id UUID NOT NULL,
    exercise_id UUID NOT NULL,
    sets_planned INT NOT NULL DEFAULT 1,
    reps_per_set_planned INT [] NOT NULL,
    sets_completed INT NOT NULL DEFAULT 0,
    reps_per_set_completed INT [] NOT NULL,
    weights_planned_lbs INT [] NOT NULL,
    weights_completed_lbs INT [] NOT NULL,
    date_completed TIMESTAMP,
    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    sort_order INT NOT NULL DEFAULT 0
);

-- +goose statementbegin
INSERT INTO workouts_exercises (
    id,
    workout_id,
    exercise_id,
    sets_planned,
    reps_per_set_planned,
    sets_completed,
    reps_per_set_completed,
    weights_planned_lbs,
    weights_completed_lbs,
    date_completed,
    updated_at,
    created_at,
    sort_order
) VALUES
-- bench press in upper body strength workout
(
    gen_random_uuid(),
    (
        SELECT id FROM workouts
        WHERE title = 'upper body strength'
    ),
    (
        SELECT id FROM exercises
        WHERE name = 'bench press'
    ),
    3,
    '{10,10,10}',
    3,
    '{10,9,8}',
    '{135,135,135}',
    '{135,135,135}',
    now() - INTERVAL '14 days',
    now(),
    now(),
    1
),

-- pull-ups in upper body strength workout
(
    gen_random_uuid(),
    (
        SELECT id FROM workouts
        WHERE title = 'upper body strength'
    ),
    (
        SELECT id FROM exercises
        WHERE name = 'pull-up'
    ),
    4,
    '{8,8,6,6}',
    4,
    '{8,7,6,5}',
    '{0,0,0,0}',
    '{0,0,0,0}',
    now() - INTERVAL '14 days',
    now(),
    now(),
    2
),

-- squats in leg day workout
(
    gen_random_uuid(),
    (
        SELECT id FROM workouts
        WHERE title = 'leg day'
    ),
    (
        SELECT id FROM exercises
        WHERE name = 'back squat'
    ),
    4,
    '{8,8,8,8}',
    4,
    '{8,8,7,6}',
    '{185,185,185,185}',
    '{185,185,185,175}',
    now() - INTERVAL '10 days',
    now(),
    now(),
    1
),

-- lunges in leg day workout
(
    gen_random_uuid(),
    (
        SELECT id FROM workouts
        WHERE title = 'leg day'
    ),
    (
        SELECT id FROM exercises
        WHERE name = 'lunge'
    ),
    3,
    '{10,10,10}',
    3,
    '{10,10,9}',
    '{0,0,0}',
    '{0,0,0}',
    now() - INTERVAL '10 days',
    now(),
    now(),
    2
);
-- +goose statementend

-- +goose Down
DROP TABLE workouts_exercises;
