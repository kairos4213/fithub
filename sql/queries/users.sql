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
) VALUES (gen_random_uuid(), now(), now(), $1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET
    hashed_password = coalesce(sqlc.narg('hashedPassword'), hashed_password),
    email = coalesce(sqlc.narg('email'), email),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreateOAuthUser :one
INSERT INTO users (
    id,
    created_at,
    updated_at,
    first_name,
    last_name,
    email,
    profile_image
) VALUES (gen_random_uuid(), now(), now(), $1, $2, $3, $4)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
