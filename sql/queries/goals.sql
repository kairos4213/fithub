-- name: CreateGoal :one
INSERT INTO goals (
  id,
  created_at,
  updated_at,
  name,
  description,
  goal_date,
  notes,
  user_id
) VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5)
RETURNING *;