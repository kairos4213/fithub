-- name: AddExercise :one
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
);
