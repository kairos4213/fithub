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

-- name: GetRandomExercisesByMuscleGroup :many
SELECT * FROM exercises
WHERE primary_muscle_group = $1
ORDER BY RANDOM()
LIMIT $2;

-- name: GetRandomExerciseExcluding :one
SELECT * FROM exercises
WHERE primary_muscle_group = $1
  AND id != ALL(@exclude_ids::uuid[])
ORDER BY RANDOM()
LIMIT 1;

-- name: GetAllMuscleGroups :many
SELECT DISTINCT muscle_group
FROM (
    SELECT primary_muscle_group AS muscle_group FROM exercises
    UNION
    SELECT secondary_muscle_group AS muscle_group FROM exercises
) AS all_groups
WHERE muscle_group IS NOT NULL
ORDER BY muscle_group;
