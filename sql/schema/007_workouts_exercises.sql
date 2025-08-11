-- +goose Up
CREATE TABLE workouts_exercises (
    id UUID PRIMARY KEY,
    workout_id UUID NOT NULL,
    exercise_id UUID NOT NULL,
    sets_planned INT NOT NULL DEFAULT 1,
    reps_per_set_planned INT [] NOT NULL,
    sets_completed INT NOT NULL DEFAULT 0,
    reps_per_set_completed INT [] NOT NULL,
    weights_planned_lbs INT [] NOT NULL,
    weights_completed_lbs INT [] NOT NULL,
    date_completed TIMESTAMP,
    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX workout_id_sort_order_idx
ON workouts_exercises (workout_id, sort_order);

-- +goose Down
DROP TABLE workouts_exercises;
