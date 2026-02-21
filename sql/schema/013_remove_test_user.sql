-- +goose Up
DELETE FROM users WHERE email = 'user@email.com';

-- +goose Down
-- No rollback for test data removal
