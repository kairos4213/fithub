-- name: CreateExercise :one
INSERT INTO exercises (
    id,
    created_at,
    updated_at,
    name,
    description,
    primary_muscle_groups,
    secondary_muscle_groups
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
    primary_muscle_groups = $3,
    secondary_muscle_groups = $4
WHERE id = $5
RETURNING *;

-- name: DeleteExercise :exec
DELETE FROM exercises
WHERE id = $1;

-- name: GetExerciseByName :many
SELECT * FROM exercises
WHERE name = $1;

-- name: GetExercisesByPrimaryMG :many
SELECT * FROM exercises
WHERE primary_muscle_groups = $1;

-- name: GetExercisesBySecondaryMG :many
SELECT * FROM exercises
WHERE secondary_muscle_groups = $1;
