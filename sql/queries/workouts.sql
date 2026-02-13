-- name: CreateWorkout :one
INSERT INTO workouts (
    id,
    created_at,
    updated_at,
    user_id,
    title,
    description,
    duration_minutes,
    planned_date
) VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetAllUserWorkouts :many
SELECT * FROM workouts
WHERE user_id = $1;

-- name: GetUpcomingUserWorkouts :many
SELECT * FROM workouts
WHERE user_id = $1 AND date_completed IS NULL
ORDER BY planned_date ASC;

-- name: GetCompletedUserWorkouts :many
SELECT * FROM workouts
WHERE user_id = $1 AND date_completed IS NOT NULL
ORDER BY date_completed DESC;

-- name: GetWorkoutByID :one
SELECT * FROM workouts
WHERE id = $1 AND user_id = $2;

-- name: UpdateWorkout :one
UPDATE workouts
SET
    updated_at = now(),
    title = $1,
    description = $2,
    duration_minutes = $3,
    planned_date = $4,
    date_completed = $5
WHERE id = $6 AND user_id = $7
RETURNING *;

-- name: DeleteWorkout :exec
DELETE FROM workouts
WHERE id = $1 AND user_id = $2;

-- name: DeleteAllUserWorkouts :exec
DELETE FROM workouts
WHERE user_id = $1;
