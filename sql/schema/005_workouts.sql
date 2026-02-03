-- +goose Up
create table workouts (
    id uuid primary key,
    user_id uuid not null references users (id) on delete cascade,
    title text not null,
    description text,
    duration_minutes integer not null,
    planned_date timestamp not null,
    date_completed timestamp,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose statementbegin
insert into workouts (
    id,
    user_id,
    title,
    description,
    duration_minutes,
    planned_date,
    date_completed,
    created_at,
    updated_at
) values
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    'upper body strength',
    'bench press, shoulder press, pull-ups, biceps curls, triceps dips.',
    60,
    now() - interval '14 days',
    now() - interval '14 days',
    now() - interval '14 days',
    now() - interval '14 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    'leg day',
    'squats, lunges, deadlifts, calf raises, hip thrusts.',
    70,
    now() - interval '10 days',
    now() - interval '10 days',
    now() - interval '10 days',
    now() - interval '10 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    'cardio intervals',
    'treadmill intervals alternating between sprint and jog.',
    45,
    now() - interval '5 days',
    null,
    now() - interval '5 days',
    now() - interval '5 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    'full body circuit',
    'mix of strength and cardio: burpees, kettlebell swings, push-ups, '
    || 'rows, mountain climbers.',
    55,
    now() + interval '2 days',
    null,
    now(),
    now()
);
-- +goose statementend

-- +goose Down
drop table workouts;
