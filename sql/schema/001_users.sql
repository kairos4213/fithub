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
  preferences JSON DEFAULT NULL
);

-- +goose Down
DROP TABLE users;