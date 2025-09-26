-- +goose Up
CREATE TABLE body_weights (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    measurement NUMERIC(5, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE muscle_masses (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    measurement NUMERIC(5, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE body_fat_percents (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    measurement NUMERIC(4, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose statementbegin
-- body_weights
INSERT INTO body_weights (
    id,
    user_id,
    measurement,
    created_at,
    updated_at
) VALUES
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    185.20,
    now() - INTERVAL '30 days',
    now() - INTERVAL '30 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    183.75,
    now() - INTERVAL '20 days',
    now() - INTERVAL '20 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    182.10,
    now() - INTERVAL '10 days',
    now() - INTERVAL '10 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    181.60,
    now(),
    now()
);

-- muscle_masses
INSERT INTO muscle_masses (
    id,
    user_id,
    measurement,
    created_at,
    updated_at
) VALUES
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    78.50,
    now() - INTERVAL '30 days',
    now() - INTERVAL '30 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    79.10,
    now() - INTERVAL '20 days',
    now() - INTERVAL '20 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    79.80,
    now() - INTERVAL '10 days',
    now() - INTERVAL '10 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    80.20,
    now(),
    now()
);

-- body_fat_percents
INSERT INTO body_fat_percents (
    id,
    user_id,
    measurement,
    created_at,
    updated_at
) VALUES
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    22.40,
    now() - INTERVAL '30 days',
    now() - INTERVAL '30 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    21.80,
    now() - INTERVAL '20 days',
    now() - INTERVAL '20 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    21.10,
    now() - INTERVAL '10 days',
    now() - INTERVAL '10 days'
),
(
    gen_random_uuid(),
    (
        SELECT id FROM users
        WHERE first_name = 'user' AND last_name = 'test'
    ),
    20.70,
    now(),
    now()
);
-- +goose statementend

-- +goose Down
DROP TABLE body_weights, muscle_masses, body_fat_percents;
