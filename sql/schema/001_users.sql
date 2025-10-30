-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(254) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    profile_image VARCHAR(255) DEFAULT NULL,
    preferences JSON DEFAULT NULL,
    is_admin BOOLEAN DEFAULT FALSE NOT NULL
);

-- +goose statementbegin
INSERT INTO users (
    id,
    created_at,
    updated_at,
    first_name,
    last_name,
    email,
    hashed_password,
    is_admin
) VALUES (
    gen_random_uuid(),
    now(),
    now(),
    'admin',
    'istrator',
    'admin@email.com',
    '$argon2id$v=19$m=65536,t=1,p=6$9JjN19qYEh68uSWUT/AVtQ$PQV3yTw5DlwbZFTaQRvOlt2QqRkTuq6OSLCrqwpqkUc',
    TRUE
);

INSERT INTO users (
    id,
    created_at,
    updated_at,
    first_name,
    last_name,
    email,
    hashed_password,
    is_admin
) VALUES (
    gen_random_uuid(),
    now(),
    now(),
    'user',
    'test',
    'user@email.com',
    '$argon2id$v=19$m=65536,t=1,p=6$9JjN19qYEh68uSWUT/AVtQ$PQV3yTw5DlwbZFTaQRvOlt2QqRkTuq6OSLCrqwpqkUc',
    FALSE
);
-- +goose statementend

-- +goose Down
DROP TABLE users;
