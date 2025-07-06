-- +goose Up
CREATE TABLE workouts_exercises (
    id UUID PRIMARY KEY,
    workout_id UUID NOT NULL,
    exercise_id UUID NOT NULL,
    sets_planned SMALLINT NOT NULL DEFAULT 1,
    reps_per_set_planned SMALLINT NOT NULL DEFAULT 1,
    sets_completed SMALLINT NOT NULL DEFAULT 0,
    reps_per_set_completed SMALLINT NOT NULL DEFAULT 0,
    weights_planned_lbs SMALLINT [] NOT NULL,
    weights_completed_lbs SMALLINT [] NOT NULL,
    date_completed TIMESTAMP,
    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE workouts_exercises;
