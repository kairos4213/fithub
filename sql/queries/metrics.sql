-- name: AddBodyWeight :one
INSERT INTO body_weights (
  id,
  created_at,
  updated_at,
  user_id,
  measurement
) VALUES ( gen_random_uuid(), NOW(), NOW(), $1, $2)
  RETURNING *;
