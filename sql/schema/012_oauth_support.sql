-- +goose Up

-- Allow OAuth-only users who have no password
ALTER TABLE users ALTER COLUMN hashed_password DROP NOT NULL;
ALTER TABLE users ALTER COLUMN hashed_password SET DEFAULT NULL;

CREATE TABLE auth_providers (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider varchar(50) NOT NULL,
    provider_user_id varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_auth_providers_user_id ON auth_providers(user_id);

-- +goose Down
DROP TABLE IF EXISTS auth_providers;
ALTER TABLE users ALTER COLUMN hashed_password SET NOT NULL;
ALTER TABLE users ALTER COLUMN hashed_password DROP DEFAULT;
