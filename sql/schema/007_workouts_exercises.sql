-- +goose Up
create table workouts_exercises (
    id uuid primary key,
    workout_id uuid not null references workouts (id) on delete cascade,
    exercise_id uuid not null references exercises (id) on delete cascade,
    sets_planned int not null default 1,
    reps_per_set_planned int [] not null,
    sets_completed int not null default 0,
    reps_per_set_completed int [] not null,
    weights_planned_lbs int [] not null,
    weights_completed_lbs int [] not null,
    date_completed timestamp,
    updated_at timestamp not null,
    created_at timestamp not null,
    sort_order int not null default 0
);

-- +goose statementbegin
insert into workouts_exercises (
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
) values (
  -- bench press in upper body strength workout
    gen_random_uuid(),
    (
        select id from workouts
        where title = 'upper body strength'
    ),
    (
        select id from exercises
        where name = 'barbell bench press'
    ),
    3,
    '{10,10,10}',
    3,
    '{10,9,8}',
    '{135,135,135}',
    '{135,135,135}',
    now() - interval '14 days',
    now(),
    now(),
    1
),

-- pull-ups in upper body strength workout
(
    gen_random_uuid(),
    (
        select id from workouts
        where title = 'upper body strength'
    ),
    (
        select id from exercises
        where name = 'pull-ups'
    ),
    4,
    '{8,8,6,6}',
    4,
    '{8,7,6,5}',
    '{0,0,0,0}',
    '{0,0,0,0}',
    now() - interval '14 days',
    now(),
    now(),
    2
),

-- squats in leg day workout
(
    gen_random_uuid(),
    (
        select id from workouts
        where title = 'leg day'
    ),
    (
        select id from exercises
        where name = 'barbell back squat'
    ),
    4,
    '{8,8,8,8}',
    4,
    '{8,8,7,6}',
    '{185,185,185,185}',
    '{185,185,185,175}',
    now() - interval '10 days',
    now(),
    now(),
    1
),

-- lunges in leg day workout
(
    gen_random_uuid(),
    (
        select id from workouts
        where title = 'leg day'
    ),
    (
        select id from exercises
        where name = 'reverse lunges'
    ),
    3,
    '{10,10,10}',
    3,
    '{10,10,9}',
    '{0,0,0}',
    '{0,0,0}',
    now() - interval '10 days',
    now(),
    now(),
    2
);
-- +goose statementend

-- +goose Down
drop table workouts_exercises;
