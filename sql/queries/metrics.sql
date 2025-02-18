-- name: AddBodyWeight :one
INSERT INTO body_weights (
  id,
  created_at,
  updated_at,
  user_id,
  measurement
) VALUES ( gen_random_uuid(), NOW(), NOW(), $1, $2)
  RETURNING *;

-- name: AddMuscleMass :one
INSERT INTO muscle_masses (
  id,
  created_at,
  updated_at,
  user_id,
  measurement
) VALUES ( gen_random_uuid(), NOW(), NOW(), $1, $2)
  RETURNING *;

-- name: AddBodyFatPerc :one
INSERT INTO body_fat_percents (
  id,
  created_at,
  updated_at,
  user_id,
  measurement
) VALUES ( gen_random_uuid(), NOW(), NOW(), $1, $2)
  RETURNING *;

-- name: GetAllBodyWeights :many
SELECT * FROM body_weights
  WHERE user_id = $1;
