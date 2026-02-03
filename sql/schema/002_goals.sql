-- +goose Up
create table goals (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    goal_name varchar(100) not null,
    description varchar(500) not null,
    goal_date timestamp not null,
    completion_date timestamp default null,
    notes text default null,
    status varchar(11) not null default 'in_progress',
    user_id uuid not null references users (id) on delete cascade,
    unique (goal_name, user_id)
);

-- +goose statementbegin
insert into goals (
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
) values
(
    gen_random_uuid(),
    now() - interval '30 days',
    now(),
    'run 5k without stopping',
    'train consistently to build endurance and complete a 5 kilometer run '
    || 'without walking breaks.',
    now() + interval '30 days',
    null,
    'currently able to run 3k comfortably. adding intervals to training.',
    'in_progress',
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    )
),
(
    gen_random_uuid(),
    now() - interval '90 days',
    now() - interval '60 days',
    'track meals daily for 30 days',
    'use a food journal or tracking app to log all meals, snacks, and drinks '
    || 'for accountability.',
    now() - interval '60 days',
    now() - interval '60 days',
    'learned a lot about portion sizes and calorie balance. built habit of '
    || 'awareness.',
    'completed',
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    )
),
(
    gen_random_uuid(),
    now() - interval '10 days',
    now(),
    'bench press 185 lbs',
    'increase strength progressively to achieve a one-rep max bench press of '
    || '185 pounds.',
    now() + interval '60 days',
    null,
    'currently at 165 lbs. adding 5 lbs per week with progressive overload.',
    'in_progress',
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    )
),
(
    gen_random_uuid(),
    now(),
    now(),
    'sleep 8 hours consistently',
    'establish a consistent sleep routine to achieve at least 8 hours of '
    || 'quality sleep per night.',
    now() + interval '90 days',
    null,
    'planning to set a bedtime alarm and reduce screen time before bed.',
    'in_progress',
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    )
);
-- +goose statementend

-- +goose Down
drop table goals;
