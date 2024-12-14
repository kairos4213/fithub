-- name: CreateUser :one
INSERT INTO users (
  id,
  created_at,
  updated_at,
  first_name,
  middle_name,
  last_name,
  email,
  hashed_password
) VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5)
RETURNING id, first_name, last_name, email;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;