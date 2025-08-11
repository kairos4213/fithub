-- name: AddExerciseToWorkout :one
INSERT INTO workouts_exercises (
    id,
    workout_id,
    exercise_id,
    sets_planned,
    reps_per_set_planned,
    sets_completed,
    reps_per_set_completed,
    weights_planned_lbs,
    weights_completed_lbs,
    updated_at,
    created_at,
    sort_order
) VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    now(),
    now(),
    (
        SELECT coalesce(max(sort_order), 0) + 1 FROM workouts_exercises
        WHERE workout_id = $1
    )
) RETURNING *;

-- name: ExercisesForWorkout :many
SELECT
    sqlc.embed(we),
    sqlc.embed(e)
FROM workouts_exercises AS we
JOIN exercises AS e
    ON we.exercise_id = e.id
WHERE we.workout_id = $1
ORDER BY we.sort_order;
