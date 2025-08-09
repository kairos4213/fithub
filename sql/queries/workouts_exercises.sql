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
    created_at
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
    now()
) RETURNING *;

-- name: ExercisesForWorkout :many
SELECT
    sqlc.embed(workouts_exercises),
    sqlc.embed(exercises)
FROM workouts_exercises
JOIN exercises
    ON workouts_exercises.exercise_id = exercises.id
WHERE workouts_exercises.workout_id = $1;
