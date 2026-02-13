-- +goose Up
ALTER TABLE goals
ADD CONSTRAINT goals_status_check
CHECK (status IN ('in_progress', 'completed'));

-- +goose Down
ALTER TABLE goals
DROP CONSTRAINT goals_status_check;
