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
