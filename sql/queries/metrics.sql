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

-- name: GetAllMuscleMasses :many
SELECT * FROM muscle_masses
  WHERE user_id = $1;

-- name: GetAllBodyFatPercs :many
SELECT * FROM body_fat_percents
  WHERE user_id = $1;

-- name: UpdateBodyWeight :one
UPDATE body_weights
  SET measurement = $1,
      updated_at = NOW()
  WHERE id = $2 AND user_id = $3 
  RETURNING *;

-- name: UpdateMuscleMass :one
UPDATE muscle_masses
  SET measurement = $1,
      updated_at = NOW()
  WHERE id = $2 AND user_id = $3
  RETURNING *;

-- name: UpdateBodyFatPerc :one
UPDATE body_fat_percents
  SET measurement = $1,
      updated_at = NOW()
  WHERE id = $2 AND user_id = $3
  RETURNING *;

-- name: DeleteBodyWeight :exec
DELETE FROM body_weights
  WHERE id = $1 AND user_id = $2;

-- name: DeleteMuscleMass :exec
DELETE FROM muscle_masses
  WHERE id = $1 AND user_id = $2;

-- name: DeleteBodyFatPerc :exec
DELETE FROM body_fat_percents
  WHERE id = $1 AND user_id = $2;
