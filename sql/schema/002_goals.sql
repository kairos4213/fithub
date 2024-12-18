-- +goose Up
CREATE TABLE goals (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR(100) NOT NULL,
  description VARCHAR(500) NOT NULL,
  goal_date TIMESTAMP NOT NULL,
  completion_date TIMESTAMP DEFAULT NULL,
  notes TEXT DEFAULT NULL,
  status VARCHAR(11) NOT NULL DEFAULT 'in_progress',
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE (name, user_id)
);

-- +goose Down
DROP TABLE goals;