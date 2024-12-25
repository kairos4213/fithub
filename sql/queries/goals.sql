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

-- name: GetAllUserGoals :many
SELECT * FROM goals
WHERE user_id = $1;

-- name: UpdateGoal :one
UPDATE goals
SET updated_at = NOW(),
    name = $1,
    description = $2,
    goal_date = $3,
    completion_date = $4,
    notes = $5,
    status = $6
WHERE id = $7
RETURNING *;