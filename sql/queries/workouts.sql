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
      NOW(),
      NOW(),
      $1,
      $2,
      $3,
      $4,
      $5
  ) RETURNING *;

-- name: GetAllUserWorkouts :many
SELECT * FROM workouts
    WHERE user_id = $1;

-- name: UpdateWorkout :one
UPDATE workouts
  SET 
    updated_at = NOW(),
    title = $1,
    description = $2,
    duration_minutes = $3,
    planned_date = $4,
    date_completed = $5
  WHERE id = $6 AND user_id = $7
  RETURNING *;
