-- +goose Up
create table body_weights (
    id uuid primary key,
    user_id uuid not null references users (id) on delete cascade,
    measurement numeric(5, 2) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create table muscle_masses (
    id uuid primary key,
    user_id uuid not null references users (id) on delete cascade,
    measurement numeric(5, 2) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create table body_fat_percents (
    id uuid primary key,
    user_id uuid not null references users (id) on delete cascade,
    measurement numeric(4, 2) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose statementbegin
-- body_weights
insert into body_weights (
    id,
    user_id,
    measurement,
    created_at,
    updated_at
) values
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    185.20,
    now() - interval '30 days',
    now() - interval '30 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    183.75,
    now() - interval '20 days',
    now() - interval '20 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    182.10,
    now() - interval '10 days',
    now() - interval '10 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    181.60,
    now(),
    now()
);

-- muscle_masses
insert into muscle_masses (
    id,
    user_id,
    measurement,
    created_at,
    updated_at
) values
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    78.50,
    now() - interval '30 days',
    now() - interval '30 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    79.10,
    now() - interval '20 days',
    now() - interval '20 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    79.80,
    now() - interval '10 days',
    now() - interval '10 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    80.20,
    now(),
    now()
);

-- body_fat_percents
insert into body_fat_percents (
    id,
    user_id,
    measurement,
    created_at,
    updated_at
) values
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    22.40,
    now() - interval '30 days',
    now() - interval '30 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    21.80,
    now() - interval '20 days',
    now() - interval '20 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    21.10,
    now() - interval '10 days',
    now() - interval '10 days'
),
(
    gen_random_uuid(),
    (
        select id from users
        where first_name = 'user' and last_name = 'test'
    ),
    20.70,
    now(),
    now()
);
-- +goose statementend

-- +goose Down
drop table body_weights, muscle_masses, body_fat_percents;
