-- name: AddBodyWeight :one
INSERT INTO body_weights (
    id,
    created_at,
    updated_at,
    user_id,
    measurement
) VALUES (gen_random_uuid(), now(), now(), $1, $2)
RETURNING *;

-- name: AddMuscleMass :one
INSERT INTO muscle_masses (
    id,
    created_at,
    updated_at,
    user_id,
    measurement
) VALUES (gen_random_uuid(), now(), now(), $1, $2)
RETURNING *;

-- name: AddBodyFatPerc :one
INSERT INTO body_fat_percents (
    id,
    created_at,
    updated_at,
    user_id,
    measurement
) VALUES (gen_random_uuid(), now(), now(), $1, $2)
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
SET
    measurement = $1,
    updated_at = now()
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: UpdateMuscleMass :one
UPDATE muscle_masses
SET
    measurement = $1,
    updated_at = now()
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: UpdateBodyFatPerc :one
UPDATE body_fat_percents
SET
    measurement = $1,
    updated_at = now()
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: DeleteBodyWeight :one
WITH deleted AS (
    DELETE FROM body_weights WHERE body_weights.id = $1 AND body_weights.user_id = $2 RETURNING body_weights.user_id
)
SELECT COUNT(*) FROM body_weights WHERE body_weights.user_id = $2;

-- name: DeleteMuscleMass :one
WITH deleted AS (
    DELETE FROM muscle_masses WHERE muscle_masses.id = $1 AND muscle_masses.user_id = $2 RETURNING muscle_masses.user_id
)
SELECT COUNT(*) FROM muscle_masses WHERE muscle_masses.user_id = $2;

-- name: DeleteBodyFatPerc :one
WITH deleted AS (
    DELETE FROM body_fat_percents WHERE body_fat_percents.id = $1 AND body_fat_percents.user_id = $2 RETURNING body_fat_percents.user_id
)
SELECT COUNT(*) FROM body_fat_percents WHERE body_fat_percents.user_id = $2;

-- name: DeleteAllBodyWeights :exec
DELETE FROM body_weights
WHERE user_id = $1;

-- name: DeleteAllMuscleMasses :exec
DELETE FROM muscle_masses
WHERE user_id = $1;

-- name: DeleteAllBodyFatPercs :exec
DELETE FROM body_fat_percents
WHERE user_id = $1;
