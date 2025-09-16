-- name: CreateExercise :one
INSERT INTO exercises (
    id,
    created_at,
    updated_at,
    name,
    description,
    primary_muscle_group,
    secondary_muscle_group
) VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: UpdateExercise :one
UPDATE exercises
SET
    updated_at = now(),
    name = $1,
    description = $2,
    primary_muscle_group = $3,
    secondary_muscle_group = $4
WHERE id = $5
RETURNING *;

-- name: DeleteExercise :exec
DELETE FROM exercises
WHERE id = $1;

-- name: GetAllExercises :many
SELECT * FROM exercises;

-- name: GetExerciseByName :one
SELECT * FROM exercises
WHERE name = $1;

-- name: GetExerciseByID :one
SELECT * FROM exercises
WHERE id = $1;

-- name: GetExerciseByKeyword :many
SELECT * FROM exercises
WHERE
    concat(name, ' ', primary_muscle_group, ' ', secondary_muscle_group)
    ILIKE '%' || sqlc.arg(word)::text || '%';

-- name: GetExercisesByPrimaryMG :many
SELECT * FROM exercises
WHERE primary_muscle_group = $1;

-- name: GetExercisesBySecondaryMG :many
SELECT * FROM exercises
WHERE secondary_muscle_group = $1;
