-- name: GetAllWorkoutTemplates :many
SELECT
    id,
    template_name,
    description,
    exercise_set_reps,
    duration_minutes
FROM workout_templates;
