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
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = COALESCE(sqlc.narg('hashedPassword'), hashed_password),
  email = COALESCE(sqlc.narg('email'), email),
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
